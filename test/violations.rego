package dev.eunomia.eks


# a pod (of a replica set, part et of a deployment) runs on a node 
runs_on[{"pod": pod, "node": node}] {
  some i, j, k
  
  input.topology[_].items[i].kind == "Pod"
  hostIP := input.topology[_].items[i].status.hostIP
  pod := input.topology[_].items[i].metadata.name
  rs := input.topology[_].items[i].metadata.ownerReferences[_].name

  input.topology[_].items[j].kind == "ReplicaSet"
  input.topology[_].items[j].metadata.name == rs
  deploy := input.topology[_].items[j].metadata.ownerReferences[_].name
  
  input.topology[_].items[k].kind == "Deployment"
  input.topology[_].items[k].metadata.name == deploy
  
  node := hosts_pod(hostIP)
}

# a node (EC2 instance) hosts a pod iff the status of the pod has 
# host IP equal to one of the private IPs of the node (EC2 instance)
hosts_pod(hostIP) = node {
  some i
  input.topology[i].type == "node"
  input.topology[i].ips[_] == hostIP
  node := input.topology[i].name
}