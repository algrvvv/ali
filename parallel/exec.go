package parallel

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"sync"

	"github.com/algrvvv/ali/utils"
)

func Exec(command Command, outputColor string, withoutOutput bool, wg *sync.WaitGroup) {
	defer wg.Done()

	commandLabel := utils.Colorize(command.Label, command.Color)
	fmt.Printf("Running command: %s\n", commandLabel)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", fmt.Sprintf("cd %s && %s", command.Path, command.Command))
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("cd %s && %s", command.Path, command.Command))
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Printf("failed to start command: [%s]\n", commandLabel)
		return
	}

	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for stdoutScanner.Scan() {
			output := utils.Colorize(stdoutScanner.Text(), outputColor)
			if !withoutOutput {
				fmt.Printf("[%s] %s\n", commandLabel, output)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for stderrScanner.Scan() {
			if !withoutOutput {
				fmt.Printf("[%s] %s\n", commandLabel, stderrScanner.Text())
			}
		}
	}()
}
