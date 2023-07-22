package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstrcutre:"DB_DRIVER"`
	DBSource      string `mapstrcutre:"DB_SOURCE"`
	ServerAddress string `mapstrcutre:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // json, xml , whatever

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
