package local

import (
	"fmt"

	"github.com/algrvvv/ali/utils"
)

func ExecuteLocal(
	command string, dir string,
	params []string,
	flags map[string]string, envs map[string]any,
	printResultCommands bool,
) error {
	cmd, err := utils.PrepareCommand(
		command,
		dir,
		params,
		flags,
		envs,
		printResultCommands,
	)
	if err != nil {
		return fmt.Errorf("failed to create command instance: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start exec command: %w", err)
	}

	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("failed to wait command: %w", err)
	}

	return nil
}
