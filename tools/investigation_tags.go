package tools

import "strings"

func buildInvestigationTags(serviceName, environment, namespace *string) map[string]string {
	tags := make(map[string]string)
	addInvestigationTag(tags, "service", serviceName)
	addInvestigationTag(tags, "environment", environment)
	addInvestigationTag(tags, "namespace", namespace)
	return tags
}

func addInvestigationTag(tags map[string]string, key string, value *string) {
	if value == nil {
		return
	}

	trimmedValue := strings.TrimSpace(*value)
	if trimmedValue == "" {
		return
	}

	tags[key] = trimmedValue
}
