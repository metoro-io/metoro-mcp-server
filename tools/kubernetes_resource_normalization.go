package tools

import (
	"strings"

	"github.com/metoro-io/metoro-mcp-server/model"
)

func normalizeMetricType(metricType model.MetricType) model.MetricType {
	if metricType == model.MetricType("kubernetes_resources") {
		return model.KubernetesResource
	}
	return metricType
}

func normalizeKubernetesAttribute(attribute string) string {
	switch strings.ToLower(attribute) {
	case "environment":
		return "Environment"
	case "namespace":
		return "Namespace"
	case "kind":
		return "Kind"
	case "resourcename", "resource_name", "resource.name":
		return "ResourceName"
	case "servicename", "service_name", "service.name":
		return "ServiceName"
	case "reason":
		return "Reason"
	case "pod_phase":
		return "pod_phase"
	default:
		return attribute
	}
}

func normalizeKubernetesAttributeMap(attributes map[string][]string) map[string][]string {
	if len(attributes) == 0 {
		return attributes
	}

	normalized := make(map[string][]string, len(attributes))
	for key, values := range attributes {
		normalizedKey := normalizeKubernetesAttribute(key)
		if isKubernetesJSONPathFilter(key) {
			normalizedKey = key
		}
		normalized[normalizedKey] = append(normalized[normalizedKey], values...)
	}
	return normalized
}

func normalizeKubernetesAttributeList(attributes []string) []string {
	if len(attributes) == 0 {
		return attributes
	}

	normalized := make([]string, len(attributes))
	for i, attribute := range attributes {
		if isKubernetesJSONPathFilter(attribute) {
			normalized[i] = attribute
			continue
		}
		normalized[i] = normalizeKubernetesAttribute(attribute)
	}
	return normalized
}

func isKubernetesJSONPathFilter(attribute string) bool {
	return strings.HasPrefix(attribute, "jsonPath:")
}
