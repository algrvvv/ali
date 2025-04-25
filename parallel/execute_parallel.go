package parallel

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/algrvvv/ali/utils"
)

func ExecuteParallel(
	entry *utils.AliasEntry, params []string,
	flags map[string]string, envs map[string]any,
	printResultCommands bool,
) {
	wg := &sync.WaitGroup{}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	if printResultCommands {
		fmt.Println("Configured commands:")
		for _, cmd := range entry.Cmds {
			fmt.Printf("[%s] -> %s\n", fmt.Sprintf("%s%s%s", utils.Colors["blue"], entry.AliasName, utils.Colors["reset"]), cmd)
		}
		fmt.Println(strings.Repeat("=", 30))
		fmt.Println()
	}

	for _, command := range entry.Cmds {
		wg.Add(1)
		go func() {
			defer wg.Done()

			cmd, err := utils.PrepareCommand(
				command,
				entry.Dir,
				params,
				flags,
				envs,
				printResultCommands,
			)
			if err != nil {
				fmt.Printf("failed to prepare command: [%s]\n", command)
				return
			}

			if err := cmd.Start(); err != nil {
				fmt.Printf("failed to start command: [%s]\n", command)
				return
			}

			if err := cmd.Wait(); err != nil {
				fmt.Printf("failed to start command: [%s]\n", command)
				return
			}
		}()
	}

	go func() {
		<-signalChan
		fmt.Println("got interrupt...")
		wg.Wait()
		os.Exit(1)
	}()

	wg.Wait()
}
