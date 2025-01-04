package common

import (
	"encoding/json"
	"fmt"
)

func PrettifyJSON(input string) (string, error) {
	// First unmarshal the input string into a generic interface{}
	var temp interface{}
	err := json.Unmarshal([]byte(input), &temp)
	if err != nil {
		return "", fmt.Errorf("invalid JSON input: %v", err)
	}

	// Marshal back to JSON with indentation
	pretty, err := json.MarshalIndent(temp, "", "    ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return string(pretty), nil
}
