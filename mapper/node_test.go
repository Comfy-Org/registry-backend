package mapper_test

import (
	"registry-backend/drip"
	"registry-backend/mapper"
	"testing"
)

// TestIsValidNodeID tests the isValidNodeID function with various inputs.
func TestIsValidNodeID(t *testing.T) {
	testCases := []struct {
		name          string
		node          *drip.Node
		expectedError string // include this field to specify what error message you expect
	}{
		{
			name:          "Valid Node ID",
			node:          &drip.Node{Id: stringPtr("validnodeid1")},
			expectedError: "",
		},
		{
			name:          "Node ID Too Long",
			node:          &drip.Node{Id: stringPtr("a12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901")},
			expectedError: "node id is too long",
		},
		{
			name:          "Invalid Node ID",
			node:          &drip.Node{Id: stringPtr("123")},
			expectedError: "invalid node id",
		},

		{
			name:          "Valid Node ID",
			node:          &drip.Node{Id: stringPtr("node1")},
			expectedError: "",
		},
		{
			name:          "Valid with dash",
			node:          &drip.Node{Id: stringPtr("node-1")},
			expectedError: "",
		},
		{
			name:          "Invalid with uppercase",
			node:          &drip.Node{Id: stringPtr("Node")},
			expectedError: "invalid node id",
		},
		{
			name:          "Invalid with special characters",
			node:          &drip.Node{Id: stringPtr("node_@")},
			expectedError: "invalid node id",
		},
		{
			name:          "Invalid start with number",
			node:          &drip.Node{Id: stringPtr("1node")},
			expectedError: "invalid node id",
		},
		{
			name:          "Invalid start with dash",
			node:          &drip.Node{Id: stringPtr("-node")},
			expectedError: "invalid node id",
		},
		{
			name:          "Empty input",
			node:          &drip.Node{Id: stringPtr("")},
			expectedError: "invalid node id",
		},
		{
			name:          "Valid all lowercase letters",
			node:          &drip.Node{Id: stringPtr("abcdefghijklmnopqrstuvwxyz")},
			expectedError: "",
		},
		{
			name:          "Valid all uppercase letters",
			node:          &drip.Node{Id: stringPtr("ABCD")},
			expectedError: "invalid node id",
		},
		{
			name:          "Valid containing underscore",
			node:          &drip.Node{Id: stringPtr("comfy_ui")},
			expectedError: "",
		},
		{
			name:          "Valid ID with hyphen",
			node:          &drip.Node{Id: stringPtr("valid-node-id")},
			expectedError: "",
		},
		{
			name:          "Valid ID with underscore",
			node:          &drip.Node{Id: stringPtr("valid_node_id")},
			expectedError: "",
		},
		{
			name:          "Valid ID with dot",
			node:          &drip.Node{Id: stringPtr("valid.node.id")},
			expectedError: "",
		},
		{
			name:          "Invalid ID with number first",
			node:          &drip.Node{Id: stringPtr("1invalidnodeid")},
			expectedError: "invalid node id",
		},
		{
			name:          "Invalid ID with consecutive dots",
			node:          &drip.Node{Id: stringPtr("invalid..nodeid")},
			expectedError: "invalid node id",
		},
		{
			name:          "Invalid ID with special character first",
			node:          &drip.Node{Id: stringPtr("-invalidnodeid")},
			expectedError: "invalid node id",
		},
		{
			name:          "Valid complex ID",
			node:          &drip.Node{Id: stringPtr("valid-node.id_1")},
			expectedError: "",
		},
		{
			name:          "Invalid ID with special characters only",
			node:          &drip.Node{Id: stringPtr("$$$$")},
			expectedError: "invalid node id",
		},
		{
			name:          "Invalid ID with leading dot",
			node:          &drip.Node{Id: stringPtr(".invalidnodeid")},
			expectedError: "invalid node id",
		},
		{
			name:          "Invalid ID with ending dot",
			node:          &drip.Node{Id: stringPtr("invalidnodeid.")},
			expectedError: "invalid node id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := mapper.ValidateNode(tc.node)
			if err != nil {
				if tc.expectedError == "" {
					t.Errorf("expected no error, got %v", err)
				} else if err.Error() != tc.expectedError {
					t.Errorf("expected error message %q, got %q", tc.expectedError, err.Error())
				}
			} else if tc.expectedError != "" {
				t.Errorf("expected error %q, got none", tc.expectedError)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
