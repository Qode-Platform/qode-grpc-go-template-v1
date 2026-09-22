# gRPC-Go template

Provisioned from [`Qode-Platform/fleet-template-v1`](https://github.com/Qode-Platform/fleet-template-v1) - the fleet
lifecycle contract with a gRPC-Go starter on top.

## Verified

Built and tested locally on Go 1.23.4 (toolchain auto-upgraded to 1.25):
`go build ./...` and `go test ./...` both pass.

## Fleet lifecycle

| step | command |
|---|---|
| install | `go mod download` |
| build | `go build -o ./.bin/app ./cmd/app` |
| start | `env PORT="$PORT" BASE_PATH="$BASE_PATH" ./.bin/app` |

## Notes

- NOT AN HTTP SERVICE: gRPC speaks its own protocol, so HEALTH_PATH is empty and fleet-proxy cannot browse it. The port listens, so bin/run's port check passes; readiness is the standard grpc.health.v1 service.
- BASE_PATH does not apply - a gRPC method is /package.Service/Method, not a URL an ingress prefix prepends to.

## Rule: everything under BASE_PATH

**BASE_PATH does not apply to this template.** The rule the HTTP templates
follow - mount every route under the `/direct/<agent>:<port>` prefix the fleet
forwards unchanged - has no counterpart here, because gRPC has no URL paths for
an ingress prefix to prepend to. A gRPC method is addressed as
`/package.Service/Method`; that string is part of the wire protocol, fixed by
the `.proto` service and package names, and neither the server nor a client may
rewrite it. Prefixing it would simply make every call unroutable.

So: ignore `BASE_PATH` here, do not try to "mount" services under it, and route
to this server by host and port (`PORT`) rather than by path. Readiness is the
standard `grpc.health.v1` service, not an HTTP health URL.

The rule comes back the moment you add an HTTP surface alongside the gRPC
server - a gRPC-gateway, a metrics endpoint, a debug page. Anything served over
HTTP must live under the prefix: add a `basePath()` helper (see the gin, echo,
fiber, chi or buffalo template) that reads `BASE_PATH` and normalises it to `""`
or `/leading/no-trailing-slash`, mount every handler under it, and prefix any
absolute URL you hand a client.
