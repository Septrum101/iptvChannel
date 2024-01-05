package config

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	confOnce sync.Once
	conf     = new(Config)
)

func ReadConfig() *Config {
	confOnce.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/" + appName)
		viper.AddConfigPath("$HOME/." + appName)
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			log.Panicf("fatal error config file: %v", err)
		}

		if err := viper.Unmarshal(conf); err != nil {
			log.Panic(err)
		}
	})

	return conf
}
