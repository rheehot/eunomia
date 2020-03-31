package dev.eunomia.eks


# a deploylment runs on a node (via its pods) 
# in the context of the namespace of the deployment
k8s_runs_on[{deploy,pod,node,hostIP}] {
  some i, j, k
  input.topology[_].items[i].kind == "Deployment"
  deploy := input.topology[_].items[i].metadata.name
  input.topology[_].items[j].kind == "Pod"
  pod := input.topology[_].items[j].metadata.name
  hostIP := input.topology[_].items[j].status.hostIP
  input.topology[k].type == "node"
  node := input.topology[k].name
  input.topology[k].ips[_] == hostIP
}

# checks if the deploy is the pod's supervisor,
# that is, there exists an ownership relation (via replica set)
k8s_supervisor [{deploy,pod}] {
  some i, j, k
  input.topology[_].items[i].kind == "Pod"
  pod := input.topology[_].items[i].metadata.name
  rs := input.topology[_].items[i].metadata.ownerReferences[_].name
  input.topology[_].items[j].kind == "ReplicaSet"
  input.topology[_].items[j].metadata.name == rs
  deploy := input.topology[_].items[j].metadata.ownerReferences[_].name
  input.topology[_].items[k].kind == "Deployment"
  input.topology[_].items[k].metadata.name == deploy
}

# a node is a member of a node group
# which is a member of a cluster
ec2_member[{cluster,node}] {
  some i, j, k 
  node := input.topology[i].name
  ng := input.topology[j].name
  cluster := input.topology[k].name
  input.topology[i].type == "node"
  input.topology[j].type == "nodegroup"
  input.topology[k].type == "cluster"
  input.topology[i].memberof == ng
  input.topology[j].memberof == cluster
}