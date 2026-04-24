package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// DownloadLinuxCommand 下载或更新 linux-command 仓库
func DownloadLinuxCommand(targetPath string) error {
	repoURL := "https://github.com/yaogh99123/linux-command"
	tmpDir := filepath.Join(os.TempDir(), "mdcli-linux-command-tmp")

	// 清理之前的临时目录
	os.RemoveAll(tmpDir)

	fmt.Printf("正在从 %s 克隆仓库...\n", repoURL)
	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone 失败: %w", err)
	}

	// 确保目标父目录存在
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 提取 command 和 dist
	dirsToMove := []string{"command", "dist"}
	for _, dir := range dirsToMove {
		src := filepath.Join(tmpDir, dir)
		dst := filepath.Join(targetPath, dir)

		// 如果目标目录已存在，先删除
		os.RemoveAll(dst)

		fmt.Printf("正在提取 %s 到 %s...\n", dir, dst)
		if err := os.Rename(src, dst); err != nil {
			// 如果 Rename 失败（可能是跨文件系统），尝试复制
			fmt.Printf("Rename 失败，尝试复制目录: %s\n", dir)
			cpCmd := exec.Command("cp", "-R", src, dst)
			if err := cpCmd.Run(); err != nil {
				return fmt.Errorf("复制目录 %s 失败: %w", dir, err)
			}
		}
	}

	// 清理临时目录
	os.RemoveAll(tmpDir)
	return nil
}
