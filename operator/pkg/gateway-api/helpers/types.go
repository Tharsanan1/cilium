// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package helpers

import (
	"fmt"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	corev1 "k8s.io/api/core/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	mcsapiv1alpha1 "sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"
)

const (
	kindGateway       = "Gateway"
	kindService       = "Service"
	kindServiceImport = "ServiceImport"
	kindSecret        = "Secret"
	kindHTTPRoute     = "HTTPRoute"
)

func IsGateway(parent gatewayv1.ParentReference) bool {
	return (parent.Kind == nil || *parent.Kind == kindGateway) && (parent.Group == nil || *parent.Group == gatewayv1.GroupName)
}

func IsTargetRefGateway(target gwapiv1a2.PolicyTargetReferenceWithSectionName) bool {
	return string(target.Kind) == kindGateway && string(target.Group) == gatewayv1.GroupName
}

func IsTargetRefHTTPRoute(target gwapiv1a2.PolicyTargetReferenceWithSectionName) bool {
	return string(target.Kind) == kindHTTPRoute && string(target.Group) == gatewayv1.GroupName
}

func GetAuthClusterName(sp *ciliumv2.SecurityPolicy) string {
	return fmt.Sprintf("auth_%s_%s", sp.GetName(), sp.GetNamespace())
}

func GetNamespacedAuthClusterName(sp *ciliumv2.SecurityPolicy, name string, namespace string) string {
	return fmt.Sprintf("%s/%sauth_%s_%s", name, namespace, sp.GetName(), sp.GetNamespace())
}

func IsService(be gatewayv1.BackendObjectReference) bool {
	return (be.Kind == nil || *be.Kind == kindService) && (be.Group == nil || *be.Group == corev1.GroupName)
}
func IsServiceImport(be gatewayv1.BackendObjectReference) bool {
	return be.Kind != nil && *be.Kind == kindServiceImport && be.Group != nil && *be.Group == mcsapiv1alpha1.GroupName
}

func IsSecret(secret gatewayv1.SecretObjectReference) bool {
	return (secret.Kind == nil || *secret.Kind == kindSecret) && (secret.Group == nil || *secret.Group == corev1.GroupName)
}
