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

```
$ opa eval \
      --input input.json \
      --data violations.rego \
      --package dev.eunomia.eks \
      --format pretty 'uses'
[
  {
    "namespace": "appmesh-system",
    "owner": "appmesh-controller",
    "pod": "appmesh-controller-6774f54fcc-d7rnd",
    "serviceaccount": "appmesh-controller",
    "type": "deployment"
  },
  {
    "namespace": "appmesh-system",
    "owner": "appmesh-inject",
    "pod": "appmesh-inject-7989448cdc-pjjmp",
    "serviceaccount": "appmesh-inject",
    "type": "deployment"
  },
  {
    "namespace": "kube-system",
    "owner": "coredns",
    "pod": "coredns-6658f9f447-kpq9s",
    "serviceaccount": "coredns",
    "type": "deployment"
  },
  {
    "namespace": "kube-system",
    "owner": "coredns",
    "pod": "coredns-6658f9f447-lfkn5",
    "serviceaccount": "coredns",
    "type": "deployment"
  },
  {
    "namespace": "kube-system",
    "owner": "aws-node",
    "pod": "aws-node-27phs",
    "serviceaccount": "aws-node",
    "type": "daemonset"
  },
  {
    "namespace": "kube-system",
    "owner": "aws-node",
    "pod": "aws-node-64x5t",
    "serviceaccount": "aws-node",
    "type": "daemonset"
  },
  {
    "namespace": "kube-system",
    "owner": "kube-proxy",
    "pod": "kube-proxy-hjpzd",
    "serviceaccount": "kube-proxy",
    "type": "daemonset"
  },
  {
    "namespace": "kube-system",
    "owner": "kube-proxy",
    "pod": "kube-proxy-tpt2p",
    "serviceaccount": "kube-proxy",
    "type": "daemonset"
  },
  {
    "namespace": "sts",
    "owner": "mehdb",
    "pod": "mehdb-0",
    "serviceaccount": "mehdb",
    "type": "statefulset"
  }
]
```