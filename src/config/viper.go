package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Init is called first to read all settings
func Init(path string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(path)

	viper.SetEnvPrefix("us")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error on init config: %s", err))
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("error on unmarshal config: %s", err))
	}
}
