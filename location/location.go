package location

import (
	"fmt"

	"k8s.io/client-go/util/jsonpath"
)

type SegmentType string

const SEGMENT_TYPE_KEY = "key"
const SEGMENT_TYPE_LIST_ITEM = "list-item"

type Location []Segment
type Segment struct {
	Type SegmentType
	Key  string
}

// FromJsonPath returns a Location based on the provided jsonpath expression.
// The jsonpath expression must be wrapped in curly braces {}.
// Only a subset of jsonpath is allowed.
func FromJsonPath(path string) (Location, error) {
	p, err := jsonpath.Parse("myparser", path)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}
	var out Location
	if err := walk(p.Root, &out); err != nil {
		return nil, fmt.Errorf("walking: %w", err)
	}
	return out, nil
}

func walk(root *jsonpath.ListNode, out *Location) error {
	for _, node := range root.Nodes {
		switch v := node.(type) {
		case *jsonpath.ListNode:
			if err := walk(v, out); err != nil {
				return err
			}
		case *jsonpath.FieldNode:
			segment := Segment{Type: SEGMENT_TYPE_KEY, Key: v.Value}
			*out = append(*out, segment)
		case *jsonpath.FilterNode:
			if len(v.Left.Nodes) == 0 {
				continue
			}
			if v.Left.Nodes[0].Type() != jsonpath.NodeField {
				continue
			}
			vv, ok := v.Left.Nodes[0].(*jsonpath.FieldNode)
			if !ok {
				return fmt.Errorf("unexpected type: %T, expected *jsonpath.FieldNode", v.Left.Nodes[0])
			}
			segment := Segment{Type: SEGMENT_TYPE_LIST_ITEM, Key: vv.Value}
			*out = append(*out, segment)
		case *jsonpath.ArrayNode:
			segment := Segment{Type: SEGMENT_TYPE_LIST_ITEM, Key: "*"}
			*out = append(*out, segment)
		default:
			return fmt.Errorf("disallowed: %s", node.String())
		}
	}
	return nil
}
