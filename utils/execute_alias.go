package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/algrvvv/ali/logger"
)

func ExecuteAlias(command string, args []string) {
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
