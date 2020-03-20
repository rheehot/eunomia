package violations.eunomia.testing

rbac_cluster_role[msg] {
  some i
  rbackind := input.rbac[i].kind
  rbackname := input.rbac[i].metadata.name
  rbackind == "ClusterRole"
  msg := sprintf("Found cluster role '%v'", [rbackname])
}

runs_on[violation] {
  some i, j
  input.topology[i].type == "nodes"
  input.topology[j].type == "pod"
  input.topology[j].hostrole == input.topology[i].role
  violation :=  sprintf("%v.%v runs on %v", [input.topology[j].name,input.topology[j].namespace, input.topology[i].name])
}