package config

import (
	"github.com/spf13/viper"
)

func InitConfig(path, name, ext string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(ext)
	err := viper.ReadInConfig()
	return err
}
