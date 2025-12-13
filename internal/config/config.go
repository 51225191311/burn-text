package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config 结构体映射config.yaml中的配置
type Config struct {
	Server    ServerConfig    ` mapstructure:"server"`
	Redis     RedisConfig     `mapstructure:"redis"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type RedisConfig struct {
	Addr       string `mapstructure:"addr"`
	Password   string `mapstructure:"password"`
	DB         int    `mapstructure:"db"`
	TTLMinutes int    `mapstructure:"ttl_minutes"`
}

type RateLimitConfig struct {
	Count         int64 `mapstructure:"count"`
	WindowSeconds int   `mapstructure:"window_seconds"`
}

var GlobalConfig *Config

func InitConfig() {
	env := os.Getenv("APP_ENV") //判断环境变量
	configName := "config"      //默认读取comfig.yaml

	if env == "Docker" {
		configName = "config_docker" //如果是Docker环境，读取config_docker.yaml
	}

	viper.SetConfigName(configName) // 配置文件名
	viper.SetConfigType("yaml")     // 配置
	viper.AddConfigPath(".")        // 搜索路径：当前目录

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
}
