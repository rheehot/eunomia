package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/open-policy-agent/opa/rego"

	"github.com/mhausenblas/kubecuddler"
)

// ClusterState describes the EKS cluster topology and RBAC settings
type ClusterState struct {
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
	rules := flag.String("rules", "", "The file path to the Rego rules to apply")
	outputcs := flag.Bool("dump-state", false, "If set, output the cluster state as a JSON blob")

	flag.Parse()
	if *cluster == "" {
		log.Fatalln("Need a cluster to operate on to continue, sorry :(")
	}
	log.Printf("Checking cluster [%v] in [%v] using rules from [%v]\n", *cluster, *region, *rules)

	var clusterstateraw string
	// let's check if we have JSON cluster state via stdin:
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			clusterstateraw += scanner.Text() + "\n"
		}
		if scanner.Err() != nil {
			log.Fatalln(fmt.Sprintf("Can't read cluster state from stdin: %v", scanner.Err().Error()))
		}
	}
	cs := NewClusterState(*cluster)
	switch {
	// no cluster state provided via stdin, need to establish the
	// state by querying the EKS, EC2 and Kubernetes API:
	case clusterstateraw == "":
		log.Printf("Establishing cluster state by querying AWS and Kubernetes APIs, this may take some time")
		// add node groups and nodes:
		err := addMNG(*region, *cluster, &cs)
		if err != nil {
			log.Fatalln(fmt.Sprintf("Can't query EKS and EC2 APIs: %v", err))
		}
		// add deployments and pods:
		err = addDeploys(&cs)
		if err != nil {
			log.Fatalln(fmt.Sprintf("Can't query Kubernetes: %v", err))
		}
	// add RBAC settings:
	// cluster state provided via stdin, so parse it:
	default:
		log.Printf("Reading cluster state from stdin")
		err := json.Unmarshal([]byte(clusterstateraw), &cs)
		if err != nil {
			log.Fatalln(fmt.Sprintf("Can't parse cluster state input: %v", err))
		}
	}

	// output cluster state if asked
	if *outputcs {
		bc, err := json.MarshalIndent(cs, "", " ")
		if err != nil {
			log.Fatalln(fmt.Sprintf("Can't serialize input: %v", err))
		}
		fmt.Println(string(bc))
	}

	// apply rules
	module, err := ioutil.ReadFile(*rules)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Can't read rules: %v", err))
	}
	query, err := rego.New(
		rego.Query("node_pods = data.dev.eunomia.eks.runs_on"),
		rego.Module("violations.rego", string(module)),
	).PrepareForEval(context.Background())
	if err != nil {
		log.Fatalln(fmt.Sprintf("Can't parse rules and/or data: %v", err))
	}
	results, err := query.Eval(context.Background(), rego.EvalInput(cs))
	if err != nil {
		log.Fatalln(fmt.Sprintf("Can't evaluate rules: %v", err))
	}
	fmt.Printf("%+v", results[0].Bindings["node_pods"])
}

// NewClusterState creates an input from scratch, setting only the cluster name
func NewClusterState(cluster string) ClusterState {
	cs := ClusterState{
		Topology: make([]interface{}, 1),
		RBAC:     make([]interface{}, 1),
	}
	cs.Topology[0] = struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}{
		"cluster",
		cluster,
	}
	return cs
}

// addMNG lists all managed nodegroups of a cluster in a region and
// adds them as well as their nodes (EC2 instances) to the topology
func addMNG(region, cluster string, cs *ClusterState) error {
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
		cs.Topology = append(cs.Topology, tmpng)
		nodes, ips, err := listNodes(region, *mnginfo.Nodegroup.Resources.AutoScalingGroups[0].Name)
		if err != nil {
			return fmt.Errorf("Can't list nodes of managed node group %v: %v", *ng, err)
		}
		for _, node := range nodes {
			tmpnode := struct {
				Type       string   `json:"type"`
				Membership string   `json:"memberof"`
				Name       string   `json:"name"`
				IPs        []string `json:"ips"`
			}{
				"node",
				*ng,
				node,
				ips[node],
			}
			cs.Topology = append(cs.Topology, tmpnode)
		}
	}
	return nil
}

// listNodes returns the EC2 instance IDs and private IP addresses
// of the nodes of an AutoScalingGroup
func listNodes(region, asg string) ([]string, map[string][]string, error) {
	nodes := []string{}
	var ips map[string][]string
	ips = make(map[string][]string)
	svc := autoscaling.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	res, err := svc.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(asg),
		},
	})
	if err != nil {
		return nodes, ips, err
	}
	for _, instance := range res.AutoScalingGroups[0].Instances {
		nodes = append(nodes, *instance.InstanceId)
		pIPs, err := listPrivateIPs(region, *instance.InstanceId)
		if err != nil {
			return nodes, ips, err
		}
		ips[*instance.InstanceId] = pIPs
	}
	return nodes, ips, nil
}

// listPrivateIPs returns the private IP addresses of an EC2 instance
func listPrivateIPs(region, ec2instanceid string) ([]string, error) {
	ips := []string{}
	svc := ec2.New(session.New(&aws.Config{
		Region: aws.String(region),
	}))
	res, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []*string{aws.String(ec2instanceid)},
			},
		},
	})
	if err != nil {
		return ips, err
	}
	for _, reservation := range res.Reservations {
		for _, instance := range reservation.Instances {
			ips = append(ips, *instance.PrivateIpAddress)
		}
	}
	return ips, nil
}

// addDeploys lists all deployments across all namespaces in the cluster
// and adds them as well as their pods to the topology of the input
func addDeploys(cs *ClusterState) error {
	res, err := kubecuddler.Kubectl(true, false, "", "get", "deploy,rs,po", "--all-namespaces", "--output", "json")
	if err != nil {
		return err
	}
	var drp interface{}
	err = json.Unmarshal([]byte(res), &drp)
	if err != nil {
		return err
	}
	cs.Topology = append(cs.Topology, drp)
	return nil
}
