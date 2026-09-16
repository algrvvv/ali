package utils

import (
	"slices"
	"strings"

	"github.com/algrvvv/ali/v2/logger"
)

func ParseUnknownFlags(args []string) map[string]string {
	reservedFlags := []string{
		"-D", "-debug", "--debug",
		"-print", "-print",
		"-L", "--local-config",
		"-local-config",
	}

	flags := make(map[string]string)
	for _, arg := range args {
		// NOTE: пропускаем зарезервированный ключ
		// TODO: придумать как избежать этого, если вдруг нужно будет использовать такой ключ
		if strings.Contains(arg, "=") {
			prepArg := strings.SplitN(arg, "=", 2)[0]
			if slices.Contains(reservedFlags, prepArg) {
				logger.SaveDebugf("got reserved flag: %s; skip", arg)
				continue
			}
		}

		if slices.Contains(reservedFlags, arg) {
			logger.SaveDebugf("got reserved flag: %s; skip", arg)
			continue
		}

		if strings.HasPrefix(arg, "-") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 2 {
				flags[parts[0]] = parts[1]
			} else {
				flags[parts[0]] = ""
			}
		}
	}
	return flags
}
