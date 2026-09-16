package utils

import (
	"fmt"
	"strings"
)

// MergeUnknownFlagsAndArgs функция для объединения флагов и аргументов в слайс
func MergeUnknownFlagsAndArgs(flags map[string]string, args string) []string {
	fields := strings.Fields(args) // сам удаляет пустые и лишние пробелы
	params := make([]string, 0, len(flags)+len(fields))

	// пушим флаги
	for flagName, flagValue := range flags {
		params = append(params, fmt.Sprintf("%s=%s", flagName, flagValue))
	}

	params = append(params, fields...)

	return params
}
