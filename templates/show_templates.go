package templates

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/algrvvv/ali/v2/logger"
)

func Show() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %v", err)
	}

	path := filepath.Join(home, ".ali/", DirName)
	logger.SaveDebugf("got path for templates: %q", path)

	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, entry := range dir {
		templName := GetTemplNameByFile(entry.Name())
		fmt.Printf("%d. %s\n", i+1, templName)
	}

	return nil
}
