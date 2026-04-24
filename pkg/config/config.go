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
	Lang     string    `yaml:"lang"`  // 语言设置：zh, en
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "mdcli", "config.yml")
}

// GetDefaultLinuxCommandPath 获取 linux-command 默认存储路径
func GetDefaultLinuxCommandPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "linux-command")
}

// LoadConfig 加载配置 (优先当前目录，次之家目录)
func LoadConfig() (*Config, string, error) {
	// 定义搜索路径
	home, _ := os.UserHomeDir()
	searchPaths := []string{
		"config.yml",
		".config.yml",
		filepath.Join(home, ".config", "mdcli", "config.yml"),
	}

	var data []byte
	var err error
	var foundPath string

	for _, path := range searchPaths {
		data, err = os.ReadFile(path)
		if err == nil {
			foundPath = path
			break
		}
	}

	var cfg Config
	if foundPath != "" {
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			return nil, foundPath, err
		}
	}

	// 注入默认项目 linux-command (始终排在第一位)
	defaultProject := Project{
		Name: "Linux Command",
		Path: GetDefaultLinuxCommandPath(),
	}

	// 检查是否已经存在同名或同路径项目（避免重复）
	exists := false
	for _, p := range cfg.Projects {
		if p.Path == defaultProject.Path {
			exists = true
			break
		}
	}

	if !exists {
		cfg.Projects = append([]Project{defaultProject}, cfg.Projects...)
	}

	return &cfg, foundPath, nil
}
