package parallel

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/algrvvv/ali/logger"
	"github.com/algrvvv/ali/utils"
	"github.com/spf13/viper"
)

const ParallelPrefix = "parallel"

type Command struct {
	Label   string `mapstructure:"label"`
	Color   string `mapstructure:"color"`
	Command string `mapstructure:"command"`
	Path    string `mapstructure:"path"`
}

func DoParrallel(aliasName, outputColor string, withoutOutput bool) {
	key := fmt.Sprintf("%s.%s", ParallelPrefix, aliasName)
	var commands []Command

	err := viper.UnmarshalKey(key, &commands)
	if err != nil {
		logger.SaveDebugf("failed to get parallel aliases")
		fmt.Println("failed to get parallel aliases")
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Configured commands:")
	for _, cmd := range commands {
		fmt.Printf("[%s] -> %s\n", utils.Colorize(cmd.Label, cmd.Color), cmd.Command)
	}

	fmt.Println(strings.Repeat("=", 30))
	fmt.Println()

	for _, cmd := range commands {
		wg.Add(1)
		go Exec(cmd, outputColor, withoutOutput, wg)
	}

	go func() {
		<-signalChan
		fmt.Println("got interrupt...")
		wg.Wait()
		os.Exit(1)
	}()

	wg.Wait()
}
