package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"mdcli/pkg/config"
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
		return nil, fmt.Errorf("配置文件中没有项目")
	}

	fmt.Printf("%s========================================%s\n", ColorBlue, ColorNC)
	fmt.Printf("%s      MDCLI 项目选择管理器%s\n", ColorBlue, ColorNC)
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
	fmt.Println(ColorYellow + " 常用操作：" + ColorYellow + ColorBold + "1-" + strconv.Itoa(len(cfg.Projects)) + "." + ColorYellow + "选择项目, " + ColorBold + "Enter." + ColorYellow + "进入搜索, " + ColorYellow + "q." + ColorYellow + "退出")
	fmt.Println(ColorYellow + " 快捷指令：" + ColorYellow + "[" + ColorBold + "Enter" + ColorYellow + "]查看详情, [" + ColorBold + "/" + ColorYellow + "]模糊搜索, [" + ColorBold + "Esc" + ColorYellow + "]重置查询")

	fmt.Println("")
	fmt.Println(ColorCyan + "--------------------------------------------------------" + ColorReset)
	fmt.Println("")
	fmt.Print(ColorCyan + ColorBold + "请选择项目编号 " + ColorReset + "[1-" + strconv.Itoa(len(cfg.Projects)) + "] (q 退出): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" || input == "Q" {
		os.Exit(0)
	}

	if input == "" {
		return nil, nil
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(cfg.Projects) {
		return nil, fmt.Errorf("无效的选择: %s", input)
	}

	selected := cfg.Projects[idx-1]
	fmt.Printf("\n%s已加载项目: %s%s\n", ColorBlue+ColorBold, selected.Name, ColorReset)

	return &selected, nil
}
