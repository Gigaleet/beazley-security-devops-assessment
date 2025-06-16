package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetNestedValue_Success(t *testing.T) {
	tests := []struct {
		obj      map[string]interface{}
		key      string
		expected interface{}
	}{
		{
			obj: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": "d",
					},
				},
			},
			key:      "a/b/c",
			expected: "d",
		},
		{
			obj: map[string]interface{}{
				"x": map[string]interface{}{
					"y": map[string]interface{}{
						"z": "a",
					},
				},
			},
			key:      "x/y/z",
			expected: "a",
		},
	}

	for _, tc := range tests {
		got, err := GetNestedValue(tc.obj, tc.key)
		if err != nil {
			t.Fatalf("GetNestedValue(%q) returned unexpected error: %v", tc.key, err)
		}
		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("GetNestedValue(%q) = %v; want %v", tc.key, got, tc.expected)
		}
	}
}

func TestGetNestedValue_EmptyKey(t *testing.T) {
	_, err := GetNestedValue(map[string]interface{}{"a": "b"}, "")
	if err == nil || err.Error() != "key is empty" {
		t.Errorf("expected empty-key error, got %v", err)
	}
}

func TestGetNestedValue_MissingPath(t *testing.T) {
	obj := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "value",
		},
	}
	_, err := GetNestedValue(obj, "a/b/c")
	if err == nil || !strings.Contains(err.Error(), `expected map at "c"`) {
	    t.Errorf("expected type-error at \"c\", got %v", err)
	}
}

func TestGetNestedValue_WrongType(t *testing.T) {
	obj := map[string]interface{}{
		"a": "not-a-map",
	}
	_, err := GetNestedValue(obj, "a/b")
	if err == nil || !strings.Contains(err.Error(), `expected map at "b"`) {
		t.Errorf("expected type-error, got %v", err)
	}
}
