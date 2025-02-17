package components

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Path() {
	// 1. 获取当前可执行文件的绝对路径
	exePath, err := os.Executable()
	if err != nil {
		//	fmt.Println("无法获取自身路径:", err)
		return
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		//	fmt.Println("无法获取自身绝对路径:", err)
		return
	}

	// 2. 检查当前目录是否是 C:\Windows\Temp
	currentDir := filepath.Dir(exePath)
	// 使用大小写不敏感比较
	wantedDir := `C:\Windows\Temp`
	if !samePath(currentDir, wantedDir) {
		// 如果不在 C:\Windows\Temp 下，则将自身复制过去并运行
		// 先构造一个目标路径，如 C:\Windows\Temp\myapp.exe
		newPath := filepath.Join(wantedDir, filepath.Base(exePath))

		// 执行拷贝
		err := copyFile(exePath, newPath)
		if err != nil {
			//	fmt.Println("拷贝自身失败:", err)
			return
		}

		// 3. 在新路径执行该程序
		cmd := exec.Command(newPath)
		cmd.Start() // 异步启动，新程序启动后，继续往下走

		// 4. 退出当前进程
		os.Exit(0)
	}

	// 如果已经在 C:\Windows\Temp 下，则继续运行
	//fmt.Println("已经在 C:\\Windows\\Temp 下，继续运行...")

	// 在这里写你想执行的逻辑
	// ...
}

// samePath 进行大小写不敏感的路径比较
func samePath(a, b string) bool {
	return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
}

// copyFile 复制文件内容到目标文件(覆盖)
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// 确保写入磁盘
	err = out.Sync()
	return err
}
