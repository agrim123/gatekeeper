package config

import (
	"github.com/spf13/viper"
)

func Init() {
	viper.AddConfigPath("configs")
	viper.SetConfigType("json")

	viper.SetConfigName("users")

	viper.ReadInConfig()

	viper.SetConfigName("roles")
	viper.MergeInConfig()

	viper.SetConfigName("plan")
	viper.MergeInConfig()

	viper.SetConfigName("servers")
	viper.MergeInConfig()
}
