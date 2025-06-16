package main

import (
	"fmt"
	"strings"
)

// GetNestedValue traverses a nested map[string]interface{} according to a
// slash-delimited key (e.g. "a/b/c") and returns the value found at that path.
// Returns an error if the key is empty, any path segment is missing, or an
// intermediate value isn't a map.
func GetNestedValue(obj map[string]interface{}, key string) (interface{}, error) {
	if key == "" {
		return nil, fmt.Errorf("key is empty")
	}

	parts := strings.Split(key, "/")
	var current interface{} = obj

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected map at %q, got %T", part, current)
		}
		val, exists := m[part]
		if !exists {
			return nil, fmt.Errorf("key %q not found", part)
		}
		current = val
	}

	return current, nil
}
