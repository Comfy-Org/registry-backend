package mapper_test

import (
	"google.golang.org/protobuf/proto"
	"registry-backend/drip"
	"registry-backend/mapper"
	"testing"
)

// TestIsValidNodeID tests the IsValidNodeID function with various inputs.
func TestIsValidNodeID(t *testing.T) {
	invalidSequenceError := "node id can only contain ASCII letters, digits, underscores, hyphens, and periods, and must not have invalid sequences"
	startEndCharError := "node id must not start or end with an underscore, hyphen, or period"
	lengthError := "node id must be between 1 and 100 characters"

	testCases := []struct {
		name          string
		node          *drip.Node
		expectedError string // specify the expected error message
	}{
		{"Valid Node ID", &drip.Node{Id: proto.String("validnodeid1")}, ""},
		{"Node ID Too Long", &drip.Node{Id: proto.String("a1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123")}, lengthError},
		{"Empty Node ID", &drip.Node{Id: proto.String("")}, lengthError},
		{"Starts with underscore", &drip.Node{Id: proto.String("_validnode")}, startEndCharError},
		{"Ends with hyphen", &drip.Node{Id: proto.String("validnode-")}, startEndCharError},
		{"Starts with dot", &drip.Node{Id: proto.String(".validnode")}, startEndCharError},
		{"Ends with dot", &drip.Node{Id: proto.String("validnode.")}, startEndCharError},
		{"Invalid sequence: consecutive dots", &drip.Node{Id: proto.String("invalid..node")}, ""},
		{"Invalid special characters", &drip.Node{Id: proto.String("invalid@node")}, invalidSequenceError},
		{"Valid Node with hyphen", &drip.Node{Id: proto.String("valid-node")}, ""},
		{"Valid Node with underscore", &drip.Node{Id: proto.String("valid_node")}, ""},
		{"Valid Node with period", &drip.Node{Id: proto.String("valid.node")}, ""},
		{"Invalid: starts with number", &drip.Node{Id: proto.String("1invalid")}, ""},
		{"Invalid: special characters only", &drip.Node{Id: proto.String("$$$")}, invalidSequenceError},
		{"Valid mixed-case Node", &drip.Node{Id: proto.String("ValidNodeID")}, ""},
		{"Valid long Node", &drip.Node{Id: proto.String("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123-456789")}, ""},
		{"Invalid: ends with underscore", &drip.Node{Id: proto.String("validnode_")}, startEndCharError},
		{"Invalid: invalid character middle", &drip.Node{Id: proto.String("validnode@")}, invalidSequenceError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, errMsg := mapper.IsValidNodeID(*tc.node.Id)
			if valid && tc.expectedError != "" {
				t.Errorf("expected error %q, got none", tc.expectedError)
			} else if !valid && errMsg != tc.expectedError {
				t.Errorf("expected error message %q, got %q", tc.expectedError, errMsg)
			}
		})
	}
}
