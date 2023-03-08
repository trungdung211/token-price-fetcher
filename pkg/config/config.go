package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func InitConfigs(configPath, configName string) {
	if configName == "" {
		configName = "config"
	}

	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName(configName)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("No valid config file is provided: %s", err.Error()))
	}
}
