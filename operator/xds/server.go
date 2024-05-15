package xds
import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	c_discovery "github.com/cilium/proxy/go/envoy/service/discovery/v3"
	envoy_server "github.com/cilium/proxy/go/pkg/server/v3"
	envoy_test "github.com/cilium/proxy/go/pkg/test/v3"
	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
)

var (
	gatewayRLeventChannel *(chan *GatewayRLEvent)
)

func init() {
	gatewayRLeventChannelObj := (make(chan *GatewayRLEvent, 10))
	gatewayRLeventChannel = &gatewayRLeventChannelObj
}

type GatewayRLEvent struct {
	Name string
	Namespace string
	BackendTrafficPolicies ciliumv2.BackendTrafficPolicyList
}

// RunServer starts an xDS server at the given port.
func RunServer() {
	logger := Logger{}
	cb := &envoy_test.Callbacks{Debug: logger.Debug}
	rateLimiterCache := GetRateLimiterCache()
	srv := envoy_server.NewServer(context.Background(), rateLimiterCache, cb)

	port := 18000
	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	// envoy_discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)
	c_discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)

	log.Printf("Management server listening on %d\n", port)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("Error while starting grpc server. Error %+v", err)
		}
	}()
	rlsPolicyCache.SetEmptySnapshot("node1")
	go func() {
		for event := range *gatewayRLeventChannel {
			ProcessEvent(event)
		}
	}()

}

func GetGatewayRLChannel() *chan *GatewayRLEvent {
	return gatewayRLeventChannel
}