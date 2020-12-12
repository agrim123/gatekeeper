package config

import (
	"github.com/spf13/viper"
)

func Init() {
	viper.AddConfigPath("configs")
	viper.SetConfigType("toml")

	viper.SetConfigName("config")
	viper.ReadInConfig()

	viper.SetConfigType("json")

	viper.SetConfigName("users")
	viper.MergeInConfig()

	viper.SetConfigName("groups")
	viper.MergeInConfig()

	viper.SetConfigName("plan")
	viper.MergeInConfig()

	viper.SetConfigName("servers")
	viper.MergeInConfig()
}
