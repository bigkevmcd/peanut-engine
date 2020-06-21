# peanut-engine

A simple service for deploying from a Git Repository using the gitops-engine.

This includes automated git synchronisation, and parsing of Kustomize files.


## Running

`peanut-engine` accepts a large number of options, but for most purposes you
only need three.

```shell
$ ./peanut-engine
Error: required flag(s) "branch", "path", "repo-url" not set
Usage:
  peanut-engine [flags]

Flags:
      --as string                      Username to impersonate for the operation
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --branch string                  Branch to checkout
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --default-namespace string       The namespace that should be used if resource namespace is not specified.By default resources are installed into the same namespace where peanut-engine is installed.
  -h, --help                           help for peanut-engine
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to a kube config. Only required if out-of-cluster
  -n, --namespace string               If present, the namespace scope for this CLI request
      --namespaced                     Switches agent into namespaced mode.
      --password string                Password for basic authentication to the API server
      --path string                    Path within the Repository to deploy
      --port int                       Port number. (default 9001)
      --prune                          Enables resource pruning. (default true)
      --repo-url string                Repository to deploy
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
      --resync duration                Resync duration (default 5m0s)
      --server string                  The address and port of the Kubernetes API server
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
      --username string                Username for basic authentication to the API server

2020/06/21 09:05:12 required flag(s) "branch", "path", "repo-url" not set
```

To synchronise a Git repository hosted in GitHub to your cluster, the simplest
command-line is:

```shell
$ peanut-engine --repo-url https://github.com/bigkevmcd/go-demo.git --branch master --path examples/kustomize/overlays/dev
```
