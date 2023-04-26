# Example Dummy Provider

`provider-dummy` is a minimal [Crossplane](https://crossplane.io/) Provider
that is meant to be used as an experimentation bed while developing Crossplane
providers. It includes a [simple server implementation](cmd/server/main.go) to
be used as its external API.

`server-dummy` exposes a single `/robots` endpoint to interact with robots. It
stores the data in memory and does not persist it. See example `curl` commands
to run by running `go run cmd/server/main.go`.

## Usage

There are three ways you can use this provider:
* Embedded: By default, the server is started as a background process when the
  provider image is run. You can use that embedded server by pointing the provider
  to talk to its localhost.
* Local: Deploying the provider to a Crossplane-installed cluster and running the server
  in the same cluster or a publicly accessible cluster.
* No-op: Deploying the provider to a Crossplane control plane in no-op mode and then
  deploying the controller and server to a different cluster that may or may not
  be publicly accessible.

### Embedded

The server is started as a background process when the provider image is run so
the only configuration we need is to point the provider to use that server.

```bash
cat <<EOF | kubectl apply -f -
apiVersion: dummy.upbound.io/v1alpha1
kind: ProviderConfig
metadata:
  name: default
spec:
  endpoint: http://127.0.0.1:8080
EOF
```

### No-Op Mode

In cases where the provider needs to make API calls to endpoints that are not
available to your control plane, you can deploy the provider in no-op mode and
deploy its controller separately to a cluster that can reach the external API.

For every versioned release of the provider, there is a corresponding no-op
package published with `-noop` suffix. For example, you can install the no-op
variant of `v0.2.0` by using `v0.2.0-noop` image tag.

See [this end to end guide](https://docs.upbound.io/knowledge-base/byoc/) for
more details.

### Local Mode

In this setup, both the fully operational provider and the server are deployed
to the same Crossplane-installed cluster.

You can install the provider in your cluster by running the following command:
```bash
cat <<EOF | kubectl apply -f -
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-dummy
spec:
    package: xpkg.upbound.io/upbound/provider-dummy:v0.2.0
EOF
```

Deploy the server to `crossplane-system` namespace:
```bash
kubectl -n crossplane-system apply -f cluster/server-deployment.yaml
```

Configure the provider to talk to the server:
```bash
cat <<EOF | kubectl apply -f -
apiVersion: dummy.upbound.io/v1alpha1
kind: ProviderConfig
metadata:
  name: default
spec:
  endpoint: http://server-dummy.default.svc.cluster.local
EOF
```

## Developing

The provider image already includes the server and it's built just like any other
Crossplane provider that uses `upbound/build` submodule. To build the provider,
run the following:
```bash
make build
```

If you'd like to build & push a separate image of the server (useful for no-op
and local modes), you can use the following command:
```bash
cd cmd/server
docker build -t upbound/server-dummy:latest .
```

### Local Setup

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