package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/algrvvv/ali/v2/logger"
	"github.com/spf13/viper"
)

func Edit(name string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %v", err)
	}

	templPath := name + ".yml"
	path := filepath.Join(home, ".ali/", DirName, templPath)
	logger.SaveDebugf("got path for template: %q", path)

	editor := viper.GetString("app.editor")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run template editor: %v", err)
	}

	return nil
}
