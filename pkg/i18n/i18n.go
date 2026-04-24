package i18n

import (
	"os"
	"strings"
)

type Language string

const (
	ZH Language = "zh"
	EN Language = "en"
)

var currentLang Language = ZH

var translations = map[Language]map[string]string{
	ZH: {
		"app_title":          "      MDCLI 项目选择管理器",
		"select_project":     "请选择项目编号",
		"invalid_selection":  "无效的选择",
		"exit_hint":          "q/0 退出",
		"back_to_projects":   "返回项目选择...",
		"search_prompt":      "搜索: ",
		"no_commands_found":  "未找到匹配的命令",
		"found_commands":     "找到 %d 个命令:",
		"common_ops":         " 常用操作：",
		"shortcut_keys":      " 快捷指令：",
		"select_project_op":  "选择项目",
		"enter_search":       "进入搜索",
		"view_details":       "查看详情",
		"fuzzy_search":       "模糊搜索",
		"reset_query":        "重置查询",
		"loaded_project":     "已加载项目: ",
		"init_failed":        "初始化失败",
		"config_no_projects": "配置文件中没有项目",
		"load_from_index":    "已从索引文件加载 %d 个条目",
		"load_from_scan":     "已从目录扫描加载 %d 个 Markdown 文件",
		"load_from_embed":    "已加载 %d 个嵌入的内置命令",
		"data_not_found":     "未能在 %s 找到有效 Markdown 数据 (条目数为 0)",
		"file_prefix":        "文件: %s",
		"cmd_not_found":      "命令 '%s' 不存在",
		"read_doc_failed":    "读取文档失败: %v (路径: %s)",
	},
	EN: {
		"app_title":          "      MDCLI Project Selector",
		"select_project":     "Please select a project number",
		"invalid_selection":  "Invalid selection",
		"exit_hint":          "q/0 to quit",
		"back_to_projects":   "Returning to project selection...",
		"search_prompt":      "Search: ",
		"no_commands_found":  "No matching commands found",
		"found_commands":     "Found %d commands:",
		"common_ops":         " Common Ops: ",
		"shortcut_keys":      " Shortcuts: ",
		"select_project_op":  "select project",
		"enter_search":       "enter search",
		"view_details":       "view details",
		"fuzzy_search":       "fuzzy search",
		"reset_query":        "reset query",
		"loaded_project":     "Project loaded: ",
		"init_failed":        "Initialization failed",
		"config_no_projects": "No projects found in config",
		"load_from_index":    "Loaded %d items from index file",
		"load_from_scan":     "Loaded %d Markdown files from directory scan",
		"load_from_embed":    "Loaded %d embedded internal commands",
		"data_not_found":     "Could not find valid Markdown data in %s (0 entries)",
		"file_prefix":        "File: %s",
		"cmd_not_found":      "Command '%s' not found",
		"read_doc_failed":    "Failed to read document: %v (Path: %s)",
	},
}

func init() {
	lang := os.Getenv("LANG")
	if strings.HasPrefix(lang, "en") {
		currentLang = EN
	} else {
		currentLang = ZH
	}
}

// T 翻译函数
func T(key string) string {
	if val, ok := translations[currentLang][key]; ok {
		return val
	}
	return key
}

// SetLanguage 手动设置语言
func SetLanguage(lang string) {
	if strings.HasPrefix(lang, "en") {
		currentLang = EN
	} else {
		currentLang = ZH
	}
}
