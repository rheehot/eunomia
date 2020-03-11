# eunomia

Towards a unified cloud native policy language.

## Problem statement

When using Kubernetes on AWS, for example in the form of [Amazon EKS](https://aws.amazon.com/eks/) 
one has to deal with both AWS IAM and Kubernetes RBAC.

This work rests on the following hypothesis:

> If both IAM and RBAC can be automatically and losslessly transformed into the 
> [Open Policy Agent](https://www.openpolicyagent.org/) policy language 
> [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/), 
> it provides for a simple and rich framework to build higher-level offerings 
> based on a unified representation of both policy languages.

Such higher-level offerings could be formal correctness proofs, query and look-up
tools, or combined IAM/RBAC visualizations.

The input of the tool is an IAM policy and/or an RBAC role and the output a
Rego representation of the combined input.

## UX design

```sh
$ ucnpl --iam test/iam/p000.json --rbac test/rbac/p000.yaml
2020/03/04 15:41:37 Parsing IAM file test/iam/p001.json
2020/03/04 15:41:38 Parsing RBAC file test/rbac/p000.yaml
2020/03/04 15:41:38 Generating unified Rego model
package unimodel123

allow {
 ...
}
```

## Related

- https://www.conftest.dev/
- https://www.openpolicyagent.org/docs/latest/policy-language/
- https://github.com/open-policy-agent/opa/issues/2098
- https://github.com/mhausenblas/rbIAM