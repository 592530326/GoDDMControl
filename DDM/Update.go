package DDM

import (
	"fmt"
	"os"
	"path/filepath"
)

func F热更新_检测热更新并直接更新() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("获取当前进程路径失败:", err)
		return
	}
	// 获取可执行文件的绝对路径
	absPath, err := filepath.Abs(exePath)
	if err != nil {
		fmt.Println("获取绝对路径失败:", err)
		return
	}
	fmt.Println("当前进程路径:", absPath)
}
