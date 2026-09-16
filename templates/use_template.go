package templates

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/algrvvv/ali/v2/logger"
	"github.com/spf13/viper"
)

func Use(lviper *viper.Viper, templName string, force bool) error {
	fmt.Println("use template: ", templName)

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %v", err)
	}

	templPath := templName + ".yml"
	path := filepath.Join(home, ".ali/", DirName, templPath)
	logger.SaveDebugf("got path for template: %q", path)

	wordDir, err := os.Getwd()
	if err != nil {
		return err
	}

	localConfigFile := filepath.Join(wordDir, ".ali")
	_, err = os.Stat(localConfigFile)
	if err == nil {
		if !force {
			fmt.Println("local config file already exists")
			fmt.Println("use --force for force rewrite")
			return nil
		}

		templViper := viper.New()
		templViper.SetConfigFile(path)
		templViper.SetConfigType("yaml")

		if err := templViper.ReadInConfig(); err != nil {
			return err
		}

		for _, key := range templViper.AllKeys() {
			value := templViper.Get(key)
			lviper.Set(key, value)
		}

		if err := lviper.WriteConfig(); err != nil {
			return err
		}

		fmt.Println("local config rewrited")
		return nil
	}

	fmt.Println("init new local config by template: ", templName)

	templConfig, err := os.Open(path)
	if err != nil {
		return err
	}
	defer templConfig.Close()

	localConfig, err := os.Create(localConfigFile)
	if err != nil {
		return err
	}
	defer localConfig.Close()

	_, err = io.Copy(localConfig, templConfig)
	if err != nil {
		return err
	}

	fmt.Println("local config from template created")
	return nil
}
