# Testing

First off, you want to install the [OPA CLI tool](https://www.openpolicyagent.org/docs/latest/#running-opa).

Then you can run any policies in the `test/` directory, for example as follows:
imagine you have an [IAM policy for an S3 bucket](s3_input_000.json) and a [Rego policy checking for certain violations](s3_rules_000.rego): 

```
$ opa eval \
      --input s3_input_000.json \
      --data s3_rules_000.rego \
      --package dev.eunomia.testing \
      's3_read_write_global[msg]' 
{
  "result": [
    {
      "expressions": [
        {
          "value": "Policy 'one fun policy' allows to read/write all S3 buckets",
          "text": "s3_read_write_global[msg]",
          "location": {
            "row": 1,
            "col": 1
          }
        }
      ],
      "bindings": {
        "msg": "Policy 'one fun policy' allows to read/write all S3 buckets"
      }
    }
  ]
}
```

Or, testing EKS access regime:


```
$ opa eval --input input.json --data violations.rego --package dev.eunomia.eks --format pretty 'k8s_runs_on[{deploy, ns, node}]' 
[
  [
    "coredns",
    "kube-system",
    "worker123"
  ]
]
```