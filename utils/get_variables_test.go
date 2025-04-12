package utils_test

import (
	"testing"

	"github.com/algrvvv/ali/utils"
)

func TestGetVariables(t *testing.T) {
	vars := map[string]string{
		"EXECUTE":   "./twentyone",
		"USER_FLAG": "--user=algrvvv",
		"NAME":      "go run main.go",
		"COMMAND":   "list -L",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{input: "{{EXECUTE}} {{USER_FLAG}}", expected: "./twentyone --user=algrvvv"},
		{input: "{{EXECUTE}} run", expected: "./twentyone run"},
		{input: "go run main.go {{COMMAND}}", expected: "go run main.go list -L"},
		{input: "{{NAME}} {{COMMAND}} -f", expected: "go run main.go list -L -f"},
		{input: "{{MISSING}} run", expected: "{{MISSING}} run"},
	}

	for _, test := range tests {
		got := utils.GetVariables(test.input, vars)
		if got != test.expected {
			t.Errorf("want: %s; got: %s", test.expected, got)
		}
		t.Logf("SUCCESS! got: %s", got)
	}
}
