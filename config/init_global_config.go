package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/algrvvv/ali/v2/logger"
	"github.com/algrvvv/ali/v2/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitGlobal(rootCmd *cobra.Command) {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err != nil {
		fmt.Println("error occurred: ", err)
		os.Exit(1)
	}

	// пропускаем если setup скип
	if cmd.Name() == "setup" || cmd.Name() == "version" || cmd.Name() == "help" {
		return
	}

	home, err := os.UserHomeDir()
	utils.CheckError(err)

	path := filepath.Join(home, ".ali")
	viper.AddConfigPath(path)
	viper.SetConfigType(utils.YamlConfigurationType)
	viper.SetConfigName("config")

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		utils.CheckError(err)
	}
	logger.SaveDebugf("using config: %s", viper.ConfigFileUsed())
}
