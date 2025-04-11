package config

import (
	"fmt"

	"github.com/seth16888/wxcommon/logger"

	"github.com/spf13/viper"
)

type Bootstrap struct {
	Server   *Server                  `yaml:"server"`
	Log      *logger.LogConfig        `yaml:"log"`
	Redis    *Redis                   `yaml:"redis"`
}

type Server struct {
	Addr    string `yaml:"addr"`
	Timeout int    `yaml:"timeout"`
}

type Redis struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	Username     string `yaml:"username"`
	DB           int    `yaml:"db"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

func ReadConfigFromFile(file string) *Bootstrap {
	if file == "" {
		file = "conf.yaml"
	}
	fmt.Println("read config from file: ", file)

	viper.SetConfigFile(file)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../conf")
	viper.AddConfigPath("~")

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	confVar := &Bootstrap{}
	if err := viper.Unmarshal(confVar); err != nil {
		panic(err)
	}

	// watch
	viper.WatchConfig()

	return confVar
}
