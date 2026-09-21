// A minimal gRPC server shaped for the fleet.
//
// NOT AN HTTP SERVICE. gRPC speaks its own protocol over the port, so the
// fleet's HTTP health check (HEALTH_PATH) cannot be satisfied and fleet-proxy
// cannot browse it — the port does listen, so bin/run's port check passes.
// Readiness is the standard grpc.health.v1 service, which grpc_health_probe
// and Kubernetes' gRPC probes speak.
//
// BASE_PATH is not applicable: a gRPC method is /package.Service/Method, not a
// URL path an ingress prefix can be prepended to.
package main

import (
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func port() string {
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		return p
	}
	return "8080"
}

func newServer() *grpc.Server {
	s := grpc.NewServer()
	hs := health.NewServer()
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, hs)
	reflection.Register(s) // so grpcurl can list services without a .proto
	return s
}

func main() {
	addr := ":" + port()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("grpc-template listening on %s", addr)
	if err := newServer().Serve(lis); err != nil {
		log.Fatal(err)
	}
}
