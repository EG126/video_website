package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	MySQL struct {
		DSN string `yaml:"dsn"`
	} `yaml:"mysql"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
	Static struct {
		BaseURL   string `yaml:"base_url"`
		UploadDir string `yaml:"upload_dir"`
	} `yaml:"static"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
}

var Conf *Config

func InitConfig() {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败:v%", err)
	}
	Conf = &Config{}
	if err := yaml.Unmarshal(data, Conf); err != nil {
		log.Fatalf("解析配置文件失败：:v%", err)
	}
}
