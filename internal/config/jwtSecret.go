package config

import (
	"github.com/spf13/viper"
)

func SecretKey(viper *viper.Viper) string {

	secretkey := viper.GetString("secretkey")

	return secretkey
}
