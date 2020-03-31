package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/eks"
)

// Input is the Rego rules input, describing the
// EKS cluster topology and RBAC settings
type Input struct {
	// Topology represents an EKS cluster with
	// managed node groups & nodes as well as
	// deployments and pods in it
	Topology []interface{} `json:"topology"`
	// RBAC is the collection of all SA, roles, cluster
	// roles, and respective bindings in the EKS cluster
	RBAC []interface{} `json:"rbac"`
}

func main() {
	cluster := flag.String("cluster", "", "The name of the cluster you want to check")
	region := flag.String("region", "eu-west-1", "The name of the AWS region to use.")
	flag.Parse()
	if *cluster == "" {
		log.Fatalln("Need a cluster to operate on to continue, sorry :(")
	}
	fmt.Printf("Checking cluster [%v] in [%v]\n", *cluster, *region)
	input := NewInput(*cluster)
	addMNG(*region, *cluster, &input)
	bc, err := json.MarshalIndent(input, "", " ")
	if err != nil {
		log.Fatalln(fmt.Sprintf("Can't serialize input: %v", err))
	}
	fmt.Println(string(bc))
}

// NewInput creates an input from scratch, setting only the cluster name
func NewInput(cluster string) Input {
	input := Input{
		Topology: make([]interface{}, 1),
		RBAC:     make([]interface{}, 1),
	}
	input.Topology[0] = struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{
		"cluster",
		cluster,
	}
	return input
}

// addMNG lists all managed nodegroups of a cluster in a region and
// adds them as well as their nodes (EC2 instances) to the topology of the input
func addMNG(region, cluster string, input *Input) error {
	svc := eks.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	mngs, err := svc.ListNodegroups(&eks.ListNodegroupsInput{
		ClusterName: aws.String(cluster),
	})
	if err != nil {
		return fmt.Errorf("Can't list managed node groups of cluster: %v", err)
	}
	for _, ng := range mngs.Nodegroups {
		mnginfo, err := svc.DescribeNodegroup(&eks.DescribeNodegroupInput{
			ClusterName:   aws.String(cluster),
			NodegroupName: aws.String(*ng),
		})
		if err != nil {
			return fmt.Errorf("Can't get details of managed node group %v : %v", *ng, err)
		}
		tmpng := struct {
			Type       string `json:"type"`
			Membership string `json:"memberof"`
			Name       string `json:"name"`
		}{
			"nodegroup",
			cluster,
			*ng,
		}
		input.Topology = append(input.Topology, tmpng)
		nodes, err := listNodes(region, *mnginfo.Nodegroup.Resources.AutoScalingGroups[0].Name)
		if err != nil {
			return fmt.Errorf("Can't list nodes of managed node group %v: %v", *ng, err)
		}
		for _, node := range nodes {
			tmpnode := struct {
				Type       string `json:"type"`
				Membership string `json:"memberof"`
				Name       string `json:"name"`
			}{
				"node",
				*ng,
				node,
			}
			input.Topology = append(input.Topology, tmpnode)
		}
	}
	return nil
}

// listNodes returns the EC2 instance IDs of the nodes of an AutoScalingGroup
func listNodes(region, asg string) ([]string, error) {
	nodes := []string{}
	svc := autoscaling.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	res, err := svc.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(asg),
		},
	})
	if err != nil {
		return nodes, err
	}
	for _, instance := range res.AutoScalingGroups[0].Instances {
		nodes = append(nodes, *instance.InstanceId)
	}
	return nodes, nil
}
