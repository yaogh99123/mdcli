package app

import (
	"fmt"
	"log"
	"os"

	"mdcli/pkg/commands"
	"mdcli/pkg/config"
	"mdcli/pkg/i18n"
	"mdcli/pkg/ui"
)

// Run 运行应用程序
func Run() {
	// 加载配置
	cfg, path, err := config.LoadConfig()
	if err != nil {
		log.Printf(i18n.T("init_failed")+": %v", err)
	} else if cfg != nil {
		ui.SetGlobalStyle(cfg.Style)
		if cfg.Lang != "" {
			i18n.SetLanguage(cfg.Lang)
		}
		if path != "" {
			fmt.Printf("%s Loaded config: %s %s\n", "\033[90m", path, "\033[0m")
		}
	}

	var cm *commands.CommandManager

	// 处理命令行参数
	if len(os.Args) < 2 {
		for {
			// 如果有多个项目，进入项目选择界面
			if cfg != nil && len(cfg.Projects) > 0 {
				project, err := ui.SelectProject(cfg)
				if err != nil {
					log.Fatalf(i18n.T("init_failed")+": %v", err)
				}

				if project == nil {
					// 用户直接按回车，使用默认管理器
					cm, err = commands.NewCommandManager()
					if err != nil {
						log.Fatalf(i18n.T("init_failed")+": %v", err)
					}
				} else {
					// 使用选中的项目路径初始化管理器
					cm, err = commands.NewCommandManagerWithSource(project.Path)
					if err != nil {
						log.Fatalf(i18n.T("init_failed")+" '%s': %v", project.Name, err)
					}
				}
			} else {
				// 没有配置项目，使用默认管理器
				cm, err = commands.NewCommandManager()
				if err != nil {
					log.Fatalf(i18n.T("init_failed")+": %v", err)
				}
			}

			// 进入交互式搜索模式
			err = ui.InteractiveSearch(cm, "")
			if err != nil && err.Error() == "ESC" {
				if cfg != nil && len(cfg.Projects) > 0 {
					fmt.Println("\n" + i18n.T("back_to_projects"))
					continue
				}
			}
			break
		}
		return
	}

	// 如果有参数，暂时使用默认管理器处理（或者您可以指定通过参数选择项目）
	cm, err = commands.NewCommandManager()
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	command := os.Args[1]

	switch command {
	case "file", "f", "-f", "--file":
		if len(os.Args) < 3 {
			fmt.Println("用法: mdcli file <文件路径>")
			os.Exit(1)
		}
		path := os.Args[2]
		err := ui.ShowFile(path)
		if err != nil {
			log.Fatalf("显示文件失败: %v", err)
		}

	case "view", "v", "-v", "--view":
		if len(os.Args) < 3 {
			fmt.Println("用法: mdcli view <命令名>")
			os.Exit(1)
		}
		name := os.Args[2]
		err := ui.ShowMarkdown(cm, name)
		if err != nil {
			log.Fatalf("显示失败: %v", err)
		}

	case "raw", "r", "-r", "--raw":
		if len(os.Args) < 3 {
			fmt.Println("用法: mdcli raw <命令名>")
			os.Exit(1)
		}
		name := os.Args[2]
		err := ui.ShowRaw(cm, name)
		if err != nil {
			log.Fatalf("显示失败: %v", err)
		}

	case "stats", "st", "-t", "--stats":
		Stats(cm)

	case "list", "l", "-l", "--list":
		results := cm.Search("")
		ui.ShowList(results)

	case "help", "h", "-h", "--help":
		PrintUsage()

	case "--version", "-V":
		fmt.Println("mdcli version 2.0")

	default:
		// 默认作为搜索关键词进入交互式模式
		_ = ui.InteractiveSearch(cm, command)
	}
}

// Stats 统计信息
func Stats(cm *commands.CommandManager) {
	fmt.Printf("\n=== 命令统计 ===\n")
	fmt.Printf("总命令数: %d\n\n", len(cm.GetCommands()))
}

// PrintUsage 打印使用说明
func PrintUsage() {
	fmt.Println("mdcli - Markdown CLI Tool")
	fmt.Println("https://github.com/codetips/mdcli")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Printf("  -V --version  Displays the program version string.\n")
	fmt.Printf("  -h --help     Displays help with available flag, subcommand, and positional value parameters.\n")
	fmt.Printf("  -l --list     List all available commands.\n")
	fmt.Printf("  -v --view     View the detailed usage of a command (beautifully rendered).\n")
	fmt.Printf("  -f --file     View a local Markdown file.\n")
	fmt.Printf("  -r --raw      Show the raw Markdown content of a command.\n")
	fmt.Printf("  -t --stats    Show statistics of the command database.\n")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  mdcli [keyword]   Enter interactive search mode (default).")
	fmt.Println("  mdcli view <cmd>  Show details of a specific command.")
	fmt.Println("  mdcli file <path> View a local Markdown file.")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  MDCLI_SOURCE  Path to the md_source directory (default: ./md_source).")
	fmt.Println("  MDCLI_STYLE   Render style (auto, dark, light, dracula, etc.).")
	fmt.Println("  MDCLI_WIDTH   Word wrap width (default: 100).")
}
