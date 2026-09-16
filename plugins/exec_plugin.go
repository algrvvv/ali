package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/algrvvv/ali/v2/logger"
	"github.com/algrvvv/ali/v2/utils"
)

func Exec(name string, params []string) {
	logger.SaveDebugf("exec plugin: %s", name)

	pluginPath, err := GetPluginDirByName(name)
	if err != nil {
		utils.PrintError("failed to get plugin", err)
		return
	}

	if _, err = os.Stat(pluginPath); err != nil {
		utils.PrintError("plugin not found", err)
		return
	}

	pluginConfigPath := filepath.Join(pluginPath, ConfigName)
	if _, err = os.Stat(pluginConfigPath); err != nil {
		utils.PrintError("plugin disabled", err)
		return
	}

	v, err := GetViperForPlugin(pluginConfigPath)
	if err != nil {
		utils.PrintError("failed to read plugin config", err)
		return
	}

	execCommand := v.GetString("exec")
	if execCommand == "" {
		fmt.Println("empty exec param for plugin")
		return
	}

	for _, param := range params {
		paramParts := strings.Split(param, "=")
		if len(paramParts) == 2 {
			execCommand += fmt.Sprintf(" -%s=%q", paramParts[0], paramParts[1])
		} else {
			execCommand += fmt.Sprintf(" -%s", param)
		}
	}
	logger.SaveDebugf("result command for exec: %s", execCommand)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe", "/C", execCommand)
	case "linux", "darwin":
		cmd = exec.Command("sh", "-c", execCommand)
	default:
		fmt.Println("unsuppored os")
		return
	}

	cmd.Dir = pluginPath
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		utils.PrintError("failed to start exec plugin command: %s", err)
		return
	}

	if err = cmd.Wait(); err != nil {
		utils.PrintError("failed to wait plugin command: %s", err)
		return
	}
}
