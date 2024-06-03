package mapper_test

import (
	"registry-backend/drip"
	"registry-backend/mapper"
	"testing"
)

// TestIsValidNodeID tests the isValidNodeID function with various inputs.
func TestIsValidNodeID(t *testing.T) {
	regexErrorMessage := "Node ID can only contain lowercase letters, numbers, hyphens, underscores, and dots. Dots cannot be consecutive or be at the start or end of the id."
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
			expectedError: regexErrorMessage,
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
			expectedError: "Node ID can only contain lowercase letters",
		},
		{
			name:          "Invalid with special characters",
			node:          &drip.Node{Id: stringPtr("node_@")},
			expectedError: regexErrorMessage,
		},
		{
			name:          "Invalid start with number",
			node:          &drip.Node{Id: stringPtr("1node")},
			expectedError: regexErrorMessage,
		},
		{
			name:          "Invalid start with dash",
			node:          &drip.Node{Id: stringPtr("-node")},
			expectedError: regexErrorMessage,
		},
		{
			name:          "Empty input",
			node:          &drip.Node{Id: stringPtr("")},
			expectedError: "node id must be between 1 and 50 characters",
		},
		{
			name:          "Valid all lowercase letters",
			node:          &drip.Node{Id: stringPtr("abcdefghijklmnopqrstuvwxyz")},
			expectedError: "",
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
			expectedError: regexErrorMessage,
		},
		{
			name:          "Invalid ID with consecutive dots",
			node:          &drip.Node{Id: stringPtr("invalid..nodeid")},
			expectedError: regexErrorMessage,
		},
		{
			name:          "Invalid ID with special character first",
			node:          &drip.Node{Id: stringPtr("-invalidnodeid")},
			expectedError: regexErrorMessage,
		},
		{
			name:          "Valid complex ID",
			node:          &drip.Node{Id: stringPtr("valid-node.id_1")},
			expectedError: "",
		},
		{
			name:          "Invalid ID with special characters only",
			node:          &drip.Node{Id: stringPtr("$$$$")},
			expectedError: regexErrorMessage,
		},
		{
			name:          "Invalid ID with leading dot",
			node:          &drip.Node{Id: stringPtr(".invalidnodeid")},
			expectedError: regexErrorMessage,
		},
		{
			name:          "Invalid ID with ending dot",
			node:          &drip.Node{Id: stringPtr("invalidnodeid.")},
			expectedError: regexErrorMessage,
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
