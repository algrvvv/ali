package utils

import (
	"fmt"

	"github.com/algrvvv/ali/logger"
)

// PrintError функция, которая выводит краткую информацию
// о том, что произошла ошибка, а в дебаг логи выводит
// подробную ошибку, которая произошла
func PrintError(msg string, err error) {
	fmt.Printf("%s\n\nuse -D to see more info\n", msg)
	logger.SaveDebugf("%s: %s", msg, err)
}
