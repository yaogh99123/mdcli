package commands

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	embeddedDataJSON []byte
	embeddedCommandFS embed.FS
)

// SetEmbeddedData 设置嵌入的数据（由主包调用）
func SetEmbeddedData(data []byte, fs embed.FS) {
	embeddedDataJSON = data
	embeddedCommandFS = fs
}

const (
	commandDir = "command"
	embedPrefix = "md_source/command"
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

// NewCommandManager 初始化命令管理器（使用默认路径）
func NewCommandManager() (*CommandManager, error) {
	// 优先检测本地数据源
	localMdSource := os.Getenv("MDCLI_SOURCE")
	if localMdSource == "" {
		cwd, _ := os.Getwd()
		localMdSource = filepath.Join(cwd, "md_source")
	}

	return NewCommandManagerWithSource(localMdSource)
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
				fmt.Printf("已从索引文件加载 %d 个条目\n", len(cm.commands))
				return cm, nil
			}
		}
	}

	// 2. 尝试直接扫描目录下的所有 .md 文件 (Direct 模式)
	if sourcePath != "" {
		err := cm.scanDirectory(sourcePath)
		if err == nil && len(cm.commands) > 0 {
			fmt.Printf("已从目录扫描加载 %d 个 Markdown 文件\n", len(cm.commands))
			return cm, nil
		}
	}

	// 3. 备选：从嵌入的数据中解析 JSON
	if len(embeddedDataJSON) > 0 {
		err := json.Unmarshal(embeddedDataJSON, &cm.commands)
		if err == nil && len(cm.commands) > 0 {
			cm.localSource = "" // 标记为嵌入模式
			fmt.Printf("已加载 %d 个嵌入的内置命令\n", len(cm.commands))
			return cm, nil
		}
	}

	return nil, fmt.Errorf("未能在 %s 找到有效 Markdown 数据 (条目数为 0)", sourcePath)
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
			// 避免重复或特殊文件
			if name == "" || strings.HasPrefix(name, ".") {
				return nil
			}
			cm.commands[name] = Command{
				Name:        name,
				Path:        path,
				Description: fmt.Sprintf("文件: %s", path),
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
		return "", fmt.Errorf("命令 '%s' 不存在", name)
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
	} else {
		// 3. 从嵌入的文件系统读取
		mdPath = filepath.Join(embedPrefix, fmt.Sprintf("%s.md", name))
		content, err = embeddedCommandFS.ReadFile(mdPath)
	}

	if err != nil {
		return "", fmt.Errorf("读取文档失败: %v (路径: %s)", err, mdPath)
	}

	return cm.formatContent(cmd, content), nil
}

// formatContent 格式化内容，添加标题和描述
func (cm *CommandManager) formatContent(cmd Command, content []byte) string {
	header := fmt.Sprintf("# %s\n\n> %s\n\n", cmd.Name, cmd.Description)
	return header + string(content)
}
