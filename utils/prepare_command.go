package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/algrvvv/ali/logger"
	"github.com/spf13/viper"
)

func PrepareCommand(
	command string, dir string,
	args []string, flags map[string]string,
	envs map[string]any, print bool,
) (*exec.Cmd, error) {
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
			if value == "" {
				command += " " + key
			} else {
				command += " " + fmt.Sprintf("%s=%s", key, value)
			}
		}
	}

	cmdArgs := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	logger.SaveDebugf("got cmd args: %s", cmdArgs)
	if strings.TrimSpace(cmdArgs) == "" {
		fmt.Println("alias not found; use ali list")
		logger.SaveDebugf("got empty args")
		return nil, errors.New("alias not found")
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
		cmd = exec.Command("cmd.exe", "/C", resultCmd)
	case "linux", "darwin":
		cmd = exec.Command("sh", "-c", resultCmd)
	default:
		logger.SaveDebugf("Unsupported OS")
		return nil, errors.New("unsupported OS")
	}

	if dir != "" && dir != "." {
		logger.SaveDebugf("entry use dir for exec: %q", dir)
		cmd.Dir = dir
		if strings.HasPrefix(dir, "~") {
			logger.SaveDebugf("user use ~ in dir param")
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, errors.New("failed to get user home dir")
			}

			logger.SaveDebugf("got home user dir: %q", home)
			dir = strings.Replace(dir, "~", home, 1)
			logger.SaveDebugf("command dir after change ~ to home dir: %q", dir)
		}
		// resultCmd = fmt.Sprintf("cd %s && %s", dir, resultCmd)
	}
	cmd.Dir = dir

	// работаем с env
	cmdEnv := os.Environ()
	for name, value := range envs {
		n := strings.ToUpper(name)

		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%v", n, value))
		logger.SaveDebugf("set new env variable: %s=%v", n, value)
	}

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Env = cmdEnv

	return cmd, nil
}
