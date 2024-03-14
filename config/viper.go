package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitViper(fileName, filePath string) error {
	viper.SetConfigName(fileName)
	viper.AddConfigPath(filePath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return fmt.Errorf("no such config file")
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("read config error")
		}
	}
	return nil
}
