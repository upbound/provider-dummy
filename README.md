# Example Dummy Provider

`provider-dummy` is a minimal [Crossplane](https://crossplane.io/) Provider
that is meant to be used as an experimentation bed while developing Crossplane
providers. It includes a [simple server implementation](cmd/server/main.go) to
be used as its external API.

`server-dummy` exposes a single `/robots` endpoint to interact with robots. It
stores the data in memory and does not persist it. See example `curl` commands
to run by running `go run cmd/server/main.go`.

## Getting Started

```bash
# This will build two images: provider-dummy and server-dummy.
make build
```

```bash
kind create cluster --wait 5m
```

```bash
# This will deploy server-dummy and port-forward its service to localhost:8080.
# And start the provider-dummy controller locally to connect to the server.
make dev
```

```bash
# Create your first Robot!
kubectl apply -f examples/iam/robots.yaml
```

```bash
kubectl get robots
```
