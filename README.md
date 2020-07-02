# peanut-engine

A simple component for deploying from a Git Repository using `argoproj/gitops-engine`.

This includes automated git synchronisation, and parsing of Kustomize files.

This is really just a thin wrapper around gitops-engine.

## Running

`peanut-engine` accepts a large number of options, but for most purposes you only need three.

To synchronise a Git repository hosted in GitHub to your cluster, the simplest
command-line is:

```shell
$ peanut-engine --repo-url https://github.com/org/repo.git --branch main --path deploy/environments/staging
```

This will run something very similar to `kubectl apply -k deploy/environments/staging` from your repo/branch combination, every 5 minutes.

### Resync frequency

By default, your resources are applied every 5 minutes, this can be configured
via the `--resync` option, this accepts "s", "m", and "h" e.g. `3h` would cause
your cluster to be synchronised every 3 hours.

## Metrics

Prometheus metrics are exposed by default at `http://service:8080/metrics`.

## Triggering  manually

Your cluster will be synchronised with the desired frequency (see [Resync frequency](#resync-frequency) above), but you can also trigger a resync manually with curl.

```shell
$ curl -X POST http://service:8080/
```

## Disable pruning

By default, `peanut-engine` will "prune" resources that don't exist in your namespace from the data you provide.

You can disable this with `--prune=false` as a command-line option.

## Namespacing mode

You can limit the namespaces that `peanut-engine` targets, by configuring
appropriately.

Modify the command-line in th Deployment to use

`--namespaced=true`

## Command-line flags

`peanut-engine` has a number of command-flags, but most of these are to allow
configuration of the Kubernetes API client.

The following flags control the behaviour of `peanut-engine` specifically.

```
 --repo-url string                Repository to deploy e.g. https://github.com/example/example.git
 --branch string                  Branch to checkout e.g. production
 --path string                    Path within the Repository to deploy e.g. deploy
 --resync duration                Resync frequency (default 5m0s)
 --auth-token string              Authentication token to use for private repositories
 --parser string                  Which parser to use kustomize, or manifest, manifest will parse non-Kustomize configurations (default "kustomize")
 --prune                          Enables resource pruning - i.e. resources not in the set will be removed
 --default-namespace string       The namespace that should be used if resource namespace is not specified.By default resources are installed into the same namespace where peanut-engine is installed.
 --namespaced                     Switches agent into namespaced mode
```

## Testing

```shell
$ go test ./...
```
