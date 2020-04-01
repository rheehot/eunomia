# Testing

First off, you want to install the [OPA CLI tool](https://www.openpolicyagent.org/docs/latest/#running-opa).

Then you can run the example as follows: imagine you have an [EKS cluster](input.json) and a Rego policy [checking violations](violations.rego): 

```
$ opa eval \
      --input input.json \
      --data violations.rego \
      --package dev.eunomia.eks \
      --format pretty 'runs_on'
[
  {
    "node": "i-00112233445566778",
    "pod": "appmesh-controller-6774f54fcc-d7rnd"
  },
  {
    "node": "i-00998877665544332",
    "pod": "appmesh-inject-7989448cdc-pjjmp"
  },
  {
    "node": "i-00112233445566778",
    "pod": "coredns-6658f9f447-kpq9s"
  },
  {
    "node": "i-00998877665544332",
    "pod": "coredns-6658f9f447-lfkn5"
  }
]
```