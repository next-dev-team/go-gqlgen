package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type GormConfig interface {
	GetDsn() string
	GetMaxIdleConn() int
	GetMaxOpenConn() int
	GetMaxLifetimeConn() time.Duration
}

type gormConfig struct {
	DSN              string `mapstructure:"dsn"`
	MaxIdleConns     int    `mapstructure:"max-idle-connection"`
	MaxOpenConns     int    `mapstructure:"max-open-connection"`
	MaxLifetimeConns int    `mapstructure:"max-lifetime-connection"` // time is second
}

func (c *gormConfig) GetDsn() string {
	return c.DSN
}

func (c *gormConfig) GetMaxIdleConn() int {
	return c.MaxIdleConns
}

func (c *gormConfig) GetMaxOpenConn() int {
	return c.MaxOpenConns
}

func (c *gormConfig) GetMaxLifetimeConn() time.Duration {
	return time.Duration(c.MaxLifetimeConns) * time.Second
}

var gormConf *gormConfig

func GetGormConfig() GormConfig {
	if gormConf == nil {
		panic("please call gormConfig or gormConfigWithViper first")
	}
	return gormConf
}

func InitDefaultGormConfig() error {
	return InitGormConfig(false, "")
}

func InitGormConfig(cmd bool, file string) error {
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
		vp.SetConfigName("gorm.toml")
	}
	if err := vp.ReadInConfig(); err != nil {
		return err
	}
	return InitGormConfigWithViper(vp)
}

func InitGormConfigWithViper(vp *viper.Viper) error {
	gormConf = &gormConfig{}
	return vp.Unmarshal(gormConf)
}
