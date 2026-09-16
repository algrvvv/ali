package vars_test

import (
	"testing"

	"github.com/algrvvv/ali/v2/vars"
)

func TestGetVariables(t *testing.T) {
	v := map[string]string{
		"execute":   "./twentyone",
		"user_flag": "--user=algrvvv",
		"name":      "go run main.go",
		"command":   "list -L",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{input: "{{execute}} {{user_flag}}", expected: "./twentyone --user=algrvvv"},
		{input: "{{execute}} run", expected: "./twentyone run"},
		{input: "go run main.go {{command}}", expected: "go run main.go list -L"},
		{input: "{{name}} {{command}} -f", expected: "go run main.go list -L -f"},
		{input: "{{missing}} run", expected: "{{missing}} run"},
	}

	for _, test := range tests {
		got := vars.GetVariables(test.input, v)
		if got != test.expected {
			t.Errorf("ERROR: want: %s; got: %s", test.expected, got)
		} else {
			t.Logf("SUCCESS! got: %s", got)
		}
	}
}
