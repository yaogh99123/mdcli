package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"mdcli/pkg/config"
	"mdcli/pkg/i18n"
)

// 定义颜色
const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorGray   = "\033[90m"
	ColorBlue   = "\033[0;34m"
	ColorCyan   = "\033[0;36m"
	ColorNC     = "\033[0m" // No Color
)

// SelectProject 展示项目列表供用户选择（基于文本美化）
func SelectProject(cfg *config.Config) (*config.Project, error) {
	if len(cfg.Projects) == 0 {
		return nil, fmt.Errorf(i18n.T("config_no_projects"))
	}

	fmt.Printf("%s========================================%s\n", ColorBlue, ColorNC)
	fmt.Printf("%s%s%s\n", ColorBlue, i18n.T("app_title"), ColorNC)
	fmt.Printf("%s========================================%s\n", ColorBlue, ColorNC)
	fmt.Println("")
	for i, p := range cfg.Projects {
		fmt.Printf("%s [%d]%s %s%-20s%s  %s(%s)%s\n",
			ColorNC, i+1, ColorReset,
			ColorNC, p.Name, ColorReset,
			ColorGray, p.Path, ColorReset)
	}
	fmt.Println("")
	fmt.Println(ColorCyan + "--------------------------------------------------------" + ColorReset)

	fmt.Println("")
	// 插入 mdcli 相关的提示信息
	fmt.Println(ColorYellow + i18n.T("common_ops") + ColorYellow + ColorBold + "1-" + strconv.Itoa(len(cfg.Projects)) + "." + ColorYellow + i18n.T("select_project_op") + ", " + ColorBold + "Enter." + ColorYellow + i18n.T("enter_search") + ", " + ColorYellow + i18n.T("exit_hint"))
	fmt.Println(ColorYellow + i18n.T("shortcut_keys") + ColorYellow + "[" + ColorBold + "Enter" + ColorYellow + "]" + i18n.T("view_details") + ", [" + ColorBold + "/" + ColorYellow + "]" + i18n.T("fuzzy_search") + ", [" + ColorBold + "Esc" + ColorYellow + "]" + i18n.T("reset_query"))

	fmt.Println("")
	fmt.Println(ColorCyan + "--------------------------------------------------------" + ColorReset)
	fmt.Println("")
	fmt.Print(ColorCyan + ColorBold + i18n.T("select_project") + " " + ColorReset + "[1-" + strconv.Itoa(len(cfg.Projects)) + "] (" + i18n.T("exit_hint") + "): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" || input == "Q" || input == "0" {
		os.Exit(0)
	}

	if input == "" {
		return nil, nil
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(cfg.Projects) {
		return nil, fmt.Errorf("%s: %s", i18n.T("invalid_selection"), input)
	}

	selected := cfg.Projects[idx-1]
	fmt.Printf("\n%s%s%s%s\n", ColorBlue+ColorBold, i18n.T("loaded_project"), selected.Name, ColorReset)

	return &selected, nil
}
