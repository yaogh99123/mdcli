package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"mdcli/pkg/commands"
	"mdcli/pkg/i18n"
	"mdcli/pkg/utils"

	"github.com/charmbracelet/glamour"
	fzf "github.com/junegunn/fzf/src"
)

var globalStyle string

// SetGlobalStyle 设置全局渲染风格
func SetGlobalStyle(s string) {
	globalStyle = s
}

// ShowMarkdown 在终端中美化显示 Markdown (从管理器读取)
func ShowMarkdown(cm *commands.CommandManager, name string) error {
	content, err := cm.GetDetail(name)
	if err != nil {
		return err
	}
	return RenderAndShow(content)
}

// ShowFile 直接显示本地 Markdown 文件
func ShowFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf(i18n.T("init_failed")+": %v", err)
	}
	return RenderAndShow(string(data))
}

// RenderAndShow 通用的渲染并显示逻辑
func RenderAndShow(content string) error {
	// 获取风格配置：环境变量 > 配置文件 > 默认 auto
	style := os.Getenv("MDCLI_STYLE")
	if style == "" {
		style = globalStyle
	}
	if style == "" {
		style = "auto"
	}

	width := 100
	if widthStr := os.Getenv("MDCLI_WIDTH"); widthStr != "" {
		if w, err := fmt.Sscanf(widthStr, "%d", &width); err == nil && w == 1 && width > 0 {
		} else {
			width = 100
		}
	}

	var opts []glamour.TermRendererOption
	switch style {
	case "dark":
		opts = append(opts, glamour.WithStandardStyle("dark"))
	case "light":
		opts = append(opts, glamour.WithStandardStyle("light"))
	case "notty":
		opts = append(opts, glamour.WithStandardStyle("notty"))
	case "pink":
		opts = append(opts, glamour.WithStandardStyle("pink"))
	case "dracula":
		opts = append(opts, glamour.WithStandardStyle("dracula"))
	case "tokyo-night":
		// 加载自定义 JSON 样式
		stylePath := "/Users/codetips/Documents/Dev_Project/github/mdcli/pkg/style/tokyo-night.json"
		styleData, err := os.ReadFile(stylePath)
		if err == nil {
			opts = append(opts, glamour.WithStylesFromJSONBytes(styleData))
		} else {
			// 如果加载失败，降级到 auto
			opts = append(opts, glamour.WithAutoStyle())
		}
	default:
		// 检查是否是其他有效的文件路径
		if _, err := os.Stat(style); err == nil {
			styleData, err := os.ReadFile(style)
			if err == nil {
				opts = append(opts, glamour.WithStylesFromJSONBytes(styleData))
			} else {
				opts = append(opts, glamour.WithAutoStyle())
			}
		} else {
			opts = append(opts, glamour.WithAutoStyle())
		}
	}

	opts = append(opts,
		glamour.WithWordWrap(width),
		glamour.WithPreservedNewLines(),
	)

	r, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		fmt.Println(content)
		return nil
	}

	out, err := r.Render(content)
	if err != nil {
		fmt.Println(content)
		return nil
	}

	return utils.ShowWithPager(out)
}

// ShowRaw 显示原始 Markdown
func ShowRaw(cm *commands.CommandManager, name string) error {
	content, err := cm.GetDetail(name)
	if err != nil {
		return err
	}
	fmt.Println(content)
	return nil
}

// ShowList 显示命令列表
func ShowList(results []commands.Command) {
	if len(results) == 0 {
		fmt.Println(i18n.T("no_commands_found"))
		return
	}

	fmt.Printf("\n"+i18n.T("found_commands")+"\n", len(results))
	fmt.Println(strings.Repeat("=", 90))

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	for i, cmd := range results {
		desc := cmd.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		fmt.Printf("%4d. %-20s %s\n", i+1, cmd.Name, desc)
	}

	fmt.Println(strings.Repeat("=", 90))
}

// InteractiveSearch 交互式搜索
func InteractiveSearch(cm *commands.CommandManager, initialQuery string) error {
	// 预先准备数据，避免每次循环重复计算
	cmds := cm.GetCommands()
	var names []string
	for name := range cmds {
		names = append(names, name)
	}
	sort.Strings(names)

	currentQuery := initialQuery

	for {
		inputChan := make(chan string)
		go func() {
			for _, name := range names {
				cmd := cmds[name]
				inputChan <- fmt.Sprintf("%s\t%s", cmd.Name, cmd.Description)
			}
			close(inputChan)
		}()

		outputChan := make(chan string, 1)

		fzfArgs := []string{
			"--reverse",
			"--height=60%",
			"--prompt=" + i18n.T("search_prompt"),
			"--bind=esc:print(ESC)+abort",
			"--bind=ctrl-c:print(CTRL-C)+abort",
			"--delimiter=\t",
		}
		if currentQuery != "" {
			fzfArgs = append(fzfArgs, "--query", currentQuery)
		}

		options, err := fzf.ParseOptions(true, fzfArgs)
		if err != nil {
			return fmt.Errorf(i18n.T("init_failed")+": %v", err)
		}

		options.Input = inputChan
		options.Output = outputChan

		code, err := fzf.Run(options)
		if code == 130 {
			// 区分 Esc 和 Ctrl+C
			// 此时 fzf 已停止，我们直接尝试从 outputChan 读取 print 的内容
			// 使用 select 防止潜在的阻塞（虽然按理说 print 应该已经完成了）
			select {
			case out := <-outputChan:
				if out == "ESC" {
					return fmt.Errorf("ESC")
				}
				if out == "CTRL-C" {
					os.Exit(0)
				}
			default:
				// 如果没有读到标识，默认按 Ctrl+C 处理直接退出
				os.Exit(0)
			}
		}

		if err != nil {
			return fmt.Errorf(i18n.T("init_failed")+" (code %d): %v", code, err)
		}

		// 处理正常选择
		select {
		case selected := <-outputChan:
			if selected != "" {
				parts := strings.Split(selected, "\t")
				if len(parts) > 0 {
					name := strings.TrimSpace(parts[0])
					ShowMarkdown(cm, name)
					currentQuery = ""
				}
			}
		default:
			// 如果没有任何输出，可能是非预期的退出，直接结束
			return nil
		}
	}
}
