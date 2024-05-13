package xds

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	corev3 "github.com/cilium/proxy/go/envoy/config/core/v3"
	gcp_types "github.com/cilium/proxy/go/pkg/cache/types"
	envoy_cachev3 "github.com/cilium/proxy/go/pkg/cache/v3"
	gcp_cache "github.com/cilium/proxy/go/pkg/cache/v3"
	gcp_resource "github.com/cilium/proxy/go/pkg/resource/v3"
	rls_config "github.com/cilium/proxy/go/ratelimit/config/ratelimit/v3"
	"github.com/cilium/cilium/operator/pkg/model"
	"github.com/cilium/cilium/operator/pkg/gateway-api/helpers"
)

// Constants relevant to the rate limit service
const (
	RateLimiterDomain     = "Default"
	maxRandomInt      int = 999999999
)

var (
	rlsPolicyCache *rateLimitPolicyCache
	gatewayRLCache map[string]*rls_config.RateLimitConfig
)

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
		xdsCache: gcp_cache.NewSnapshotCache(false, IDHash{}, nil),
	}
	gatewayRLCache = make(map[string]*rls_config.RateLimitConfig)
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

func updateSnapshot(label string) bool {
	rlConfigs := make([]gcp_types.Resource, 0)
	for _, rls := range gatewayRLCache {
		rlConfigs = append(rlConfigs, rls)
	}

	version := fmt.Sprint(rand.Int(rand.Reader, maxRandomBigInt()))
	snap, err := gcp_cache.NewSnapshot(version, map[gcp_resource.Type][]gcp_types.Resource{
		gcp_resource.RateLimitConfigType: rlConfigs,
	})
	if err != nil {
		return false
	}
	if err := snap.Consistent(); err != nil {
		return false
	}

	if err := rlsPolicyCache.xdsCache.SetSnapshot(context.Background(), label, snap); err != nil {
		return false
	}
	return true

}

func maxRandomBigInt() *big.Int {
	return big.NewInt(int64(maxRandomInt))
}

func ProcessEvent(event *GatewayRLEvent) {
	cilium_log.Infof("Processing event: %s", event.Name)
	var rls = &rls_config.RateLimitConfig{
		Name:        fmt.Sprintf("%s-%s", event.Name, event.Namespace),
		Domain:      fmt.Sprintf("%s-%s", event.Name, event.Namespace),
		Descriptors: []*rls_config.RateLimitDescriptor{},
	}

	if len(event.BackendTrafficPolicies.Items) > 0 {
		for _, btp := range event.BackendTrafficPolicies.Items {
			if btp.Spec.RateLimit != nil && btp.Spec.RateLimit.Global != nil  {
				for index, rule := range btp.Spec.RateLimit.Global.Rules {
					var unit rls_config.RateLimitUnit
					switch rule.Limit.Unit {
					case "Second":
						unit = rls_config.RateLimitUnit_SECOND
					case "Minute":
						unit = rls_config.RateLimitUnit_MINUTE
					case "Hour":
						unit = rls_config.RateLimitUnit_HOUR
					case "Day":
						unit = rls_config.RateLimitUnit_DAY
					default:
						unit = rls_config.RateLimitUnit_UNKNOWN
					}
					ns := helpers.NamespaceDerefOr(btp.Spec.TargetRef.Namespace, btp.GetNamespace())
					if len(rule.ClientSelectors) == 0 {
						key, value := model.GetRatelimitKeyAndValueForHttprouteRL(ns, string(btp.Spec.TargetRef.Name), btp.GetName(), btp.GetNamespace(), index)
						rls.Descriptors = append(rls.Descriptors, &rls_config.RateLimitDescriptor{
							Key: key,
							Value: value,
							RateLimit: &rls_config.RateLimitPolicy{
								Unit:            unit,
								RequestsPerUnit: uint32(rule.Limit.Requests),
							},
						})
					} else {
						rateLimit := &rls_config.RateLimitPolicy{
							Unit:            unit,
							RequestsPerUnit: uint32(rule.Limit.Requests),
						}
						descriptor := &rls_config.RateLimitDescriptor{}
						rls.Descriptors = append(rls.Descriptors, descriptor)
						var lastDescriptor *rls_config.RateLimitDescriptor
						for indexCS, cs := range rule.ClientSelectors {
							if len(cs.Headers) > 0 {
								for indexHeader, header := range cs.Headers {
									key := model.GetRatelimitKeyForHttprouteHeaderRl(ns, string(btp.Spec.TargetRef.Name), btp.GetName(), btp.GetNamespace(), indexHeader, indexCS, header.Name)
									if header.Value != nil {
										descriptor.Key = key
										descriptor.Value = *header.Value
										descriptorNew :=  &rls_config.RateLimitDescriptor{}
										descriptor.Descriptors = []*rls_config.RateLimitDescriptor{
											descriptorNew,
										}
										descriptor = descriptorNew
									} else {
										descriptor.Key = key
										descriptorNew :=  &rls_config.RateLimitDescriptor{}
										descriptor.Descriptors = []*rls_config.RateLimitDescriptor{
											descriptorNew,
										}
										descriptor = descriptorNew
									}
									lastDescriptor = descriptor
								}
							}
						}
						descriptor.RateLimit = rateLimit
						// Delete last descriptor list
						lastDescriptor.Descriptors = []*rls_config.RateLimitDescriptor{}
					}
				}
			}
		}
	}
	gatewayRLCache[fmt.Sprintf("%s-%s", event.Name, event.Namespace)] = rls
	updateSnapshot("node1")
}