# eunomia

Towards a unified cloud native policy language.

## Problem statement

When using Kubernetes on AWS, for example in the form of [Amazon EKS](https://aws.amazon.com/eks/) 
one has to deal with both AWS IAM and Kubernetes RBAC.

This work rests on the following hypothesis: if both IAM and RBAC can be
automatically and losslessly transformed into the 
[Open Policy Agent](https://www.openpolicyagent.org/) policy language 
[Rego](https://www.openpolicyagent.org/docs/latest/policy-language/), 
it provides for a simple and rich framework to build higher-level offerings 
based on a unified representation of both policy languages. Such higher-level 
offerings could be formal correctness proof, query and look-up, or visualizations.

## UX design

```sh
$ ucnpl --iam test/iam/p000.json --rbac test/rbac/p000.yaml
Parsing policies
Generating unified Rego model
Result:
---
package unimodel123

allow {
 ...
}
```

## Related

- https://www.openpolicyagent.org/docs/latest/policy-language/
- https://github.com/open-policy-agent/opa/issues/2098
- https://github.com/mhausenblas/rbIAM