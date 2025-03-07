package parallel

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/algrvvv/ali/utils"
)

func Exec(command Command, outputColor string, withoutOutput bool, wg *sync.WaitGroup) {
	defer wg.Done()

	commandLabel := utils.Colorize(command.Label, command.Color)
	label := utils.Colorize(fmt.Sprintf("[%s]", command.Label), command.Color)
	fmt.Printf("Running command: %s\n", commandLabel)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", fmt.Sprintf("cd %s && %s", command.Path, command.Command))
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("cd %s && %s", command.Path, command.Command))
	}

	// елси нет четко задонного цвета логов, пытаемся сохранить исходный.
	if outputColor == "" {
		cmd.Env = append(os.Environ(), "FORCE_COLOR=1")
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
		defer func() {
			fmt.Printf("%s %s exited.\n", label, command.Command)
		}()
		for stdoutScanner.Scan() {
			var output string
			if outputColor != "" {
				output = utils.Colorize(stdoutScanner.Text(), outputColor)
			} else {
				output = stdoutScanner.Text()
			}

			if !withoutOutput {
				fmt.Printf("%s %s\n", label, output)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for stderrScanner.Scan() {
			if !withoutOutput {
				fmt.Printf("%s %s\n", utils.Colorize(fmt.Sprintf("[%s]", command.Label), command.Color), stderrScanner.Text())
			}
		}
	}()
}
