package config

import (
	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("access")  // name of config file (without extension)
	viper.SetConfigType("json")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("configs") // path to look for the config file in
	viper.ReadInConfig()           // Find and read the config file

	viper.SetConfigName("roles")
	viper.MergeInConfig()

	viper.SetConfigName("actions")
	viper.MergeInConfig()
}
