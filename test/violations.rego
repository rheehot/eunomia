package violations.eunomia.testing

# what is a cluster role
rbac_cluster_role[msg] {
  some i
  rbackind := input.rbac[i].kind
  rbackname := input.rbac[i].metadata.name
  rbackind == "ClusterRole"
  msg := sprintf("Found cluster role '%v'", [rbackname])
}