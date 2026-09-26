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
| start | `env PORT="$PORT" ./.bin/app` |

## Notes

- NOT AN HTTP SERVICE: gRPC speaks its own protocol, so HEALTH_PATH is empty and fleet-proxy cannot browse it. The port listens, so bin/run's port check passes; readiness is the standard grpc.health.v1 service.
