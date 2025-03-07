package utils

import "fmt"

var Colors = map[string]string{
	"red":     "\033[38;5;210m",
	"green":   "\033[32m",
	"yellow":  "\033[38;5;221m",
	"blue":    "\033[38;5;110m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"gray":    "\033[90m",
	"orange":  "\033[38;5;214m",
	"pink":    "\033[38;5;207m",
	"lime":    "\033[38;5;10m",
	"white":   "\033[37m",
	"reset":   "\033[0m",
}

func Colorize(text, color string) string {
	colorCode, exists := Colors[color]
	if !exists {
		colorCode = Colors["white"]
	}

	return fmt.Sprintf("%s%s%s", colorCode, text, Colors["reset"])
}
