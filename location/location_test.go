package location_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shomron/gatekeper-jsonpath-poc/location"
)

func Test_FromJsonPath(t *testing.T) {
	tests := []struct {
		input    string
		expected location.Location
		isError  bool
	}{
		{
			input: `{$.spec.containers[?(@.name)].securityContext}`,
			expected: location.Location{
				location.Segment{
					Type: location.SEGMENT_TYPE_KEY,
					Key:  "spec",
				},
				location.Segment{
					Type: location.SEGMENT_TYPE_KEY,
					Key:  "containers",
				},
				location.Segment{
					Type: location.SEGMENT_TYPE_LIST_ITEM,
					Key:  "name",
				},
				location.Segment{
					Type: location.SEGMENT_TYPE_KEY,
					Key:  "securityContext",
				},
			},
		},
		{
			input: `{$.spec.containers[?(@.name=="foo")].ports[?(@.containerPort)]}`,
			expected: location.Location{
				location.Segment{
					Type: location.SEGMENT_TYPE_KEY,
					Key:  "spec",
				},
				location.Segment{
					Type: location.SEGMENT_TYPE_KEY,
					Key:  "containers",
				},
				location.Segment{
					Type: location.SEGMENT_TYPE_LIST_ITEM,
					Key:  "name",
				},
				location.Segment{
					Type: location.SEGMENT_TYPE_KEY,
					Key:  "ports",
				},
				location.Segment{
					Type: location.SEGMENT_TYPE_LIST_ITEM,
					Key:  "containerPort",
				},
			},
		},
		{
			input: `{$.spec.containers[*]}`,
			expected: location.Location{
				location.Segment{
					Type: location.SEGMENT_TYPE_KEY,
					Key:  "spec",
				},
				location.Segment{
					Type: location.SEGMENT_TYPE_KEY,
					Key:  "containers",
				},
				location.Segment{
					Type: location.SEGMENT_TYPE_LIST_ITEM,
					Key:  "*",
				},
			},
		},
		{
			input:   `{$.spec..containers[*]}`,
			isError: true,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {

			result, err := location.FromJsonPath(tc.input)
			if (err != nil) != tc.isError {
				t.Errorf("error mismatch: got: %v, expected?: %v", err, tc.isError)
				return
			}
			diff := cmp.Diff(tc.expected, result)
			if diff != "" {
				t.Errorf("conversion failed: %v", diff)
			}
		})
	}
}
