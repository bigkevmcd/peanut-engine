# peanut-engine

A simple service for deploying from a Git Repository using `argoproj/gitops-engine`.

This includes automated git synchronisation, and parsing of Kustomize files.

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

## Testing

```shell
$ go test ./...
```
