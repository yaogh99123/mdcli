package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Project 代表一个 MD 项目
type Project struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

// Config 代表整体配置
type Config struct {
	Projects []Project `yaml:"projects"`
	Style    string    `yaml:"style"` // 渲染风格：auto, dark, light, pink, dracula 等
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "mdcli", "config.yml")
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	path := GetConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
