package xds

import (
	gcp_cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	rls_config "github.com/envoyproxy/go-control-plane/ratelimit/config/ratelimit/v3"
	gcp_resource "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	gcp_types "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"fmt"
	"crypto/rand"
	"math/big"
	"context"
)

// Constants relevant to the rate limit service
const (
	RateLimiterDomain                    = "Default"
	maxRandomInt             int    = 999999999
)
var rlsPolicyCache *rateLimitPolicyCache
// IDHash uses ID field as the node hash.
type IDHash struct{}

// ID uses the node ID field
func (IDHash) ID(node *corev3.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Id
}


func init() {
	rlsPolicyCache = &rateLimitPolicyCache{
		xdsCache:                  gcp_cache.NewSnapshotCache(false, IDHash{}, nil),
	}
}

// GetRateLimiterCache returns xds server cache for rate limiter service.
func GetRateLimiterCache() envoy_cachev3.SnapshotCache {
	return rlsPolicyCache.xdsCache
}

type rateLimitPolicyCache struct {
	// xdsCache is the snapshot cache for the rate limiter service
	xdsCache gcp_cache.SnapshotCache
}


// SetEmptySnapshot sets an empty snapshot into the apiCache for the given label
// this is used to set empty snapshot when there are no APIs available for a label
func (r *rateLimitPolicyCache) SetEmptySnapshot(label string) bool {
	var rls = &rls_config.RateLimitConfig{
		Name:        RateLimiterDomain,
		Domain:      RateLimiterDomain,
		Descriptors: []*rls_config.RateLimitDescriptor{},
	}
	version := fmt.Sprint(rand.Int(rand.Reader, maxRandomBigInt()))
	snap, err := gcp_cache.NewSnapshot(version, map[gcp_resource.Type][]gcp_types.Resource{
		gcp_resource.RateLimitConfigType: {
			rls,
		},
	})
	if err != nil {
		return false
	}
	if err := snap.Consistent(); err != nil {
		return false
	}

	if err := r.xdsCache.SetSnapshot(context.Background(), label, snap); err != nil {
		return false
	}
	return true
}

func maxRandomBigInt() *big.Int {
	return big.NewInt(int64(maxRandomInt))
}