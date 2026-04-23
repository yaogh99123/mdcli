package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ShowWithPager 使用系统分页器显示内容
func ShowWithPager(content string) error {
	// 尝试使用环境变量中的 PAGER，如果没有则使用 less
	pagerCmd := os.Getenv("PAGER")
	if pagerCmd == "" {
		pagerCmd = "less"
	}

	// 检查是否在管道中（如果是，则不使用分页器）
	stat, _ := os.Stdout.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// 输出被重定向，直接打印
		fmt.Print(content)
		return nil
	}

	var args []string
	if pagerCmd == "less" {
		args = []string{"-R", "-i"}
	} else {
		args = []string{"-R"}
	}

	cmd := exec.Command(pagerCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 创建管道
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Print(content)
		return nil
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Print(content)
		return nil
	}

	// 写入内容
	io.WriteString(stdin, content)
	stdin.Close()

	// 等待命令完成
	return cmd.Wait()
}
