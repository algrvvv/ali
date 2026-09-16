package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/algrvvv/ali/v2/logger"
	"github.com/spf13/viper"
)

func CreateNew(templName string, force bool) error {
	templPath := templName + ".yml"
	logger.SaveDebugf("init new template")
	fmt.Println("init new template by name: ", templName)

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %v", err)
	}

	path := filepath.Join(home, ".ali/", DirName, templPath)
	logger.SaveDebugf("got path for template: %q", path)

	_, err = os.Stat(path)
	if err == nil {
		if force {
			// if err := os.Remove(path); err != nil {
			// 	return fmt.Errorf("failed to remove old template: %v", err)
			// }

			// NOTE: удалять файл не обязательно, так как os.Create выглядит так:
			//  return OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)
			// как видно, функция использует O_TRUNC, а значит перезапишет файл
			fmt.Printf("old template by name %q truncated\n", templName)
		} else {
			return ErrTemplAlreadyExists
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create templ file: %v", err)
	}
	file.Close()

	fmt.Println("template file created")

	editor := viper.GetString("app.editor")
	if editor == "" {
		editor = "vi"
	}
	logger.SaveDebugf("got editor: %s", editor)

	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run template editor: %v", err)
	}

	fmt.Println()
	fmt.Println("template saved by name: ", templName)
	fmt.Printf("for run: ali templ %s\n", templName)
	fmt.Printf("for edit: ali templ %s --edit\n", templName)

	return nil
}
