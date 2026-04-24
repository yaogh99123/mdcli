package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mdcli/pkg/i18n"
)

const (
	commandDir = "command"
)

// Command 表示命令的索引信息
type Command struct {
	Name        string `json:"n"` // 命令名称
	Path        string `json:"p"` // 路径
	Description string `json:"d"` // 描述
}

// CommandManager 命令管理器
type CommandManager struct {
	commands    map[string]Command
	localSource string // 本地数据源路径（如果存在）
}


// NewCommandManagerWithSource 根据指定路径初始化命令管理器
func NewCommandManagerWithSource(sourcePath string) (*CommandManager, error) {
	cm := &CommandManager{
		commands:    make(map[string]Command),
		localSource: sourcePath,
	}

	// 1. 尝试加载 dist/data.json (Legacy 模式)
	localDataJSON := filepath.Join(sourcePath, "dist", "data.json")
	if _, err := os.Stat(localDataJSON); err == nil {
		data, err := os.ReadFile(localDataJSON)
		if err == nil {
			err = json.Unmarshal(data, &cm.commands)
			if err == nil && len(cm.commands) > 0 {
				fmt.Printf(i18n.T("load_from_index")+"\n", len(cm.commands))
				return cm, nil
			}
		}
	}

	// 2. 尝试直接扫描目录下的所有 .md 文件 (Direct 模式)
	if sourcePath != "" {
		err := cm.scanDirectory(sourcePath)
		if err == nil && len(cm.commands) > 0 {
			fmt.Printf(i18n.T("load_from_scan")+"\n", len(cm.commands))
			return cm, nil
		}
	}


	return nil, fmt.Errorf(i18n.T("data_not_found"), sourcePath)
}

// scanDirectory 扫描目录下的所有 .md 文件
func (cm *CommandManager) scanDirectory(dir string) error {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	
	err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 跳过错误的路径
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			name := strings.TrimSuffix(info.Name(), ".md")
			// 避免重复运行特殊文件
			if name == "" || strings.HasPrefix(name, ".") {
				return nil
			}
			cm.commands[name] = Command{
				Name:        name,
				Path:        path,
				Description: fmt.Sprintf(i18n.T("file_prefix"), path),
			}
		}
		return nil
	})
	
	return err
}

// Search 查询命令（支持模糊搜索）
func (cm *CommandManager) Search(keyword string) []Command {
	var results []Command

	// 这里可以保留原有的逻辑，或者直接返回 map 供 UI 侧处理
	for _, cmd := range cm.commands {
		results = append(results, cmd)
	}

	return results
}

// GetCommands 获取所有命令 map
func (cm *CommandManager) GetCommands() map[string]Command {
	return cm.commands
}

// GetDetail 获取命令详情
func (cm *CommandManager) GetDetail(name string) (string, error) {
	cmd, exists := cm.commands[name]
	if !exists {
		return "", fmt.Errorf(i18n.T("cmd_not_found"), name)
	}

	var content []byte
	var err error
	var mdPath string

	// 1. 如果 cmd.Path 是绝对路径且文件存在，直接读取 (Direct 模式)
	if filepath.IsAbs(cmd.Path) {
		if _, err := os.Stat(cmd.Path); err == nil {
			content, err = os.ReadFile(cmd.Path)
			if err == nil {
				return cm.formatContent(cmd, content), nil
			}
		}
	}

	// 2. 如果有本地源 (Legacy 模式)
	if cm.localSource != "" {
		mdPath = filepath.Join(cm.localSource, "command", fmt.Sprintf("%s.md", name))
		content, err = os.ReadFile(mdPath)
	}

	if err != nil {
		return "", fmt.Errorf(i18n.T("read_doc_failed"), err, mdPath)
	}

	return cm.formatContent(cmd, content), nil
}

// formatContent 格式化内容，添加标题和描述
func (cm *CommandManager) formatContent(cmd Command, content []byte) string {
	header := fmt.Sprintf("# %s\n\n> %s\n\n", cmd.Name, cmd.Description)
	return header + string(content)
}
