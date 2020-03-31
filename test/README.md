# Testing

First off, you want to install the [OPA CLI tool](https://www.openpolicyagent.org/docs/latest/#running-opa).

Then you can run the example as follows: imagine you have an [EKS cluster](input.json) and a Rego policy [checking violations](violations.rego): 

```
$ opa eval --input input.json --data violations.rego --package dev.eunomia.eks --format pretty 'k8s_runs_on' 
[
  [
    "coredns",
    "kube-system",
    "worker123"
  ]
]
```