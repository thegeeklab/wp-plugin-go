package cli

import (
	"encoding/json"
	"fmt"
)

// coerceToString converts any JSON value to a string representation.
func coerceToString(val any) string {
	if val == nil {
		return ""
	}

	switch t := val.(type) {
	case map[string]any, []any:
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprint(t)
		}

		return string(b)
	default:
		return fmt.Sprint(t)
	}
}

// convertNestedMap converts map[string]map[string]any to map[string]map[string]string.
func convertNestedMap(rawNested map[string]map[string]any) map[string]map[string]string {
	result := make(map[string]map[string]string, len(rawNested))

	for group, rawGroup := range rawNested {
		result[group] = make(map[string]string, len(rawGroup))

		for k, val := range rawGroup {
			result[group][k] = coerceToString(val)
		}
	}

	return result
}

// convertFlatMap converts map[string]any to map[string]string.
func convertFlatMap(raw map[string]any) map[string]string {
	result := make(map[string]string, len(raw))

	for k, val := range raw {
		result[k] = coerceToString(val)
	}

	return result
}
