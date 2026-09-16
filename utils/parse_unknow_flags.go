package utils

import (
	"fmt"
	"slices"
	"strings"

	"github.com/algrvvv/ali/v2/logger"
)

func ParseUnknownArgs(args []string) (string, map[string]string) {
	reservedFlags := []string{
		"-D", "-debug", "--debug",
		"--print", "-print",
		"-L", "--local-config",
		"-local-config",
	}

	params := strings.Builder{}
	flags := make(map[string]string)

	for _, arg := range args {
		// NOTE: в данный момент нет зарезервированных флагов с поддержкой аргументов
		//if strings.Contains(arg, "=") {
		//	prepArg := strings.SplitN(arg, "=", 2)[0]
		//	if slices.Contains(reservedFlags, prepArg) {
		//		logger.SaveDebugf("got reserved flag: %s; skip", arg)
		//		continue
		//	}
		//}

		// NOTE: пропускаем зарезервированный ключ
		if slices.Contains(reservedFlags, arg) {
			logger.SaveDebugf("got reserved flag: %s; skip", arg)
			continue
		}

		// явное прокидывание параметра
		// проверяем длину после сплита и сохраняем
		if strings.Contains(arg, "=") && strings.Contains(arg, "-") {
			parts := strings.Split(arg, "=")
			if len(parts) == 2 {
				flags[parts[0]] = parts[1]
				continue
			}
		}

		params.WriteString(
			fmt.Sprintf("%s ", arg),
		)
	}

	return strings.TrimSpace(params.String()), flags
}
