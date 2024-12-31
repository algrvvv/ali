package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/algrvvv/ali/logger"
)

func ExecuteAlias(command string, args []string, flags map[string]string) {
	// проверяем аргументы, чтобы при пробелах в них мы не получили их как разные аргументы
	for i := range args {
		if strings.Contains(args[i], " ") {
			logger.SaveDebugf(
				"founded arg that contains spaces: %s; quotation marks enabled",
				args[i],
			)
			args[i] = fmt.Sprintf("\"%s\"", args[i])
		}
	}

	for key, value := range flags {
		k := fmt.Sprintf("<%s>", strings.ReplaceAll(key, "-", ""))
		logger.SaveDebugf("parse command for find flag: %s with value: %s", k, value)
		if strings.Contains(command, k) {
			command = strings.ReplaceAll(command, k, value)
		} else {
			command += " " + key
		}
	}

	cmdArgs := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	logger.SaveDebugf("got cmd args: %s", cmdArgs)
	if strings.TrimSpace(cmdArgs) == "" {
		fmt.Println("alias not found; use ali list")
		logger.SaveDebugf("got empty args")
		return
	}

	cmd := exec.Command("sh", "-c", cmdArgs)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	CheckError(err)

	err = cmd.Wait()
	CheckError(err)
}
