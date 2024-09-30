package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config 结构体存储所有配置信息
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig 存储服务器相关配置
type ServerConfig struct {
	Port string `yaml:"port"`
}

// DatabaseConfig 存储数据库连接相关配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

// LoadConfig 从指定的文件路径加载配置
func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
