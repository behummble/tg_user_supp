package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Redis RedisConfig `yaml:"redis"`
	Bot BotConfig `yaml:"bot"`
}

type RedisConfig struct {
	Host string `yaml:"host" env:"DB_HOST" env-default:"127.0.0.1"`
	Port string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
}

type BotConfig struct {
	Token string `yaml:"token" env:"BOT_TOKEN"`
	UpdateTimeout int `yaml:"timeout" env-default:"10"`
}

func MustLoad() *Config {
	path := loadPath()
	if path == "" {
		panic("Can`t read config file")
	}

	return loadConfig(path)
}

func loadPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()
	if path == "" {
		curDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if runtime.GOOS == "windows" {
			path = fmt.Sprintf("%s/config/config.yaml", curDir)
		} else {
			path = fmt.Sprintf("%s\\config\\config.yaml", curDir)
		}
	}

	return path
}

func loadConfig(path string) *Config {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	
	return &cfg
}