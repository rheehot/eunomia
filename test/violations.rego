package dev.eunomia.eks

# a node is a member of a node group
# which is a member of a cluster
ec2_member[{cluster, node}] {
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

# a deploylment runs on a node (via its pods) 
# in the context of the namespace of the deployment
k8s_runs_on[{deploy, ns, node}] {
  some i, j, k
  pod := input.topology[j].name
  ns := input.topology[j].namespace
  deploy := input.topology[k].name
  node := input.topology[i].name
  input.topology[i].type == "node"
  input.topology[j].type == "pod"
  input.topology[k].type == "deployment"
  input.topology[j].hostrole == input.topology[i].role
  # requires that there's a supervisor relation 
  # between the deploy and the pod:
  k8s_supervisor
  ec2_member
}

# checks if the deploy is the pod's supervisor,
# that is, there exists an ownership relation
k8s_supervisor = {deploy, pod} {
 some i, j
 deploy := input.topology[j].name
 pod := input.topology[i].name
 input.topology[i].type == "pod"
 input.topology[j].type == "deployment"
 input.topology[i].owner == input.topology[j].name
 input.topology[i].namespace = input.topology[j].namespace
}

# checks if something is a cluster role
# k8s_cluster_role[{name}] {
#   some i
#   kind := input.rbac[i].kind
#   name := input.rbac[i].metadata.name
#   kind == "ClusterRole"
# }