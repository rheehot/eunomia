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

## Usage

```sh
$ cat test/input.json | go run . --cluster appmeshga --rules test/violations.rego
2020/04/01 11:14:57 Checking cluster [appmeshga] in [eu-west-1] using rules from [test/violations.rego]
2020/04/01 11:14:57 Reading cluster state from stdin
[map[node:i-00112233445566778 pod:appmesh-controller-6774f54fcc-d7rnd] map[node:i-00998877665544332 pod:appmesh-inject-7989448cdc-pjjmp] map[node:i-00112233445566778 pod:coredns-6658f9f447-kpq9s] map[node:i-00998877665544332 pod:coredns-6658f9f447-lfkn5]]
```

## Related

- https://www.conftest.dev/
- https://www.openpolicyagent.org/docs/latest/policy-language/
- https://github.com/open-policy-agent/opa/issues/2098
- https://github.com/mhausenblas/rbIAM