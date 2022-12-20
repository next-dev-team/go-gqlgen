package config

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type ServerConfig interface {
	GetLogLevel() string
	GetPort() string
}

type serverConfig struct {
	LogLevel string `mapstructure:"log-level"`
	Port     string `mapstructure:"port"`
}

var config *serverConfig

func GetServerConfig() ServerConfig {
	if config == nil {
		panic("please call GetServerConfig or GetServerConfigWithViper first")
	}
	return config
}

func (c *serverConfig) GetLogLevel() string {
	return c.LogLevel
}

func (c *serverConfig) GetPort() string {
	return c.Port
}

func InitDefaultServerConfig() error {
	return InitServerConfig(false, "")
}

func InitServerConfig(cmd bool, file string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	vp := viper.New()
	vp.SetConfigType("toml")
	if cmd {
		vp.SetConfigFile(file)
	} else {
		vp.AddConfigPath(filepath.Join(wd, "configs"))
		vp.SetConfigName("server.toml")
	}
	if err := vp.ReadInConfig(); err != nil {
		log.Err(err).Msg("read server config error")
		return err
	}
	return InitServerConfigWithViper(vp)
}

func InitServerConfigWithViper(vp *viper.Viper) error {
	config = &serverConfig{}
	return vp.Unmarshal(config)
}
