package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/algrvvv/ali/logger"
	"github.com/spf13/viper"
)

func ExecuteAlias(command string, args []string, flags map[string]string, print bool) {
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
		logger.SaveDebugf("got key: %s", key)

		preparedKey := strings.TrimLeft(key, "-")
		if print && preparedKey == "print" {
			logger.SaveDebugf("user want to print result")
			continue
		}

		if strings.Contains(key, "V_") {
			varToChange := strings.Replace(key, "V_", "", 1)
			varToChange = strings.TrimLeft(varToChange, "-")
			logger.SaveDebugf("key: %s - contains V; var to change: %s", key, varToChange)

			varKey := fmt.Sprintf("vars.%s", varToChange)
			logger.SaveDebugf("got var key for change: %s", varKey)
			viper.Set(varKey, value)
			continue
		}

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

	resultCmd := cmdArgs
	vars, err := GetVars()
	if err != nil {
		logger.SaveDebugf("failed to get all vars: %v", err)
		fmt.Println("failed to get vars. skip")
	} else {
		logger.SaveDebugf("got vars: %v", vars)
		resultCmd = GetVariables(cmdArgs, vars)
	}

	logger.SaveDebugf("result command to execute: %s", resultCmd)
	if print {
		fmt.Println("command: ", resultCmd)
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe", "/C", cmdArgs)
	case "linux", "darwin":
		cmd = exec.Command("sh", "-c", cmdArgs)
	default:
		logger.SaveDebugf("Unsupported OS")
		return
	}

	// cmd := exec.Command("sh", "-c", cmdArgs)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	CheckError(err)

	err = cmd.Wait()
	CheckError(err)
}
