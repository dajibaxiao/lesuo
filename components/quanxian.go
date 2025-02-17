package components

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

// 这里我们需要自己定义 TokenElevation 为 20。
// 参考：
// https://learn.microsoft.com/en-us/windows/win32/api/winnt/ne-winnt-token_information_class
const (
	TOKEN_ELEVATION = 20
	SW_SHOWNORMAL   = 1
)

// 声明要用到的 Windows API。
// 1. ShellExecuteW 来执行提权操作
// 2. 其余操作用 syscall 包封装
var (
	shell32           = syscall.NewLazyDLL("shell32.dll")
	procShellExecuteW = shell32.NewProc("ShellExecuteW")
)

// IsAdmin 判断当前进程是否在管理员权限下运行
func IsAdmin() bool {
	// 1. 获取当前进程
	process, err := syscall.GetCurrentProcess()
	if err != nil {
		return false
	}

	// 2. 打开进程令牌 (TOKEN_QUERY)
	var token syscall.Token
	err = syscall.OpenProcessToken(process, syscall.TOKEN_QUERY, &token)
	if err != nil {
		return false
	}
	defer token.Close()

	// 3. 调用 GetTokenInformation(TokenElevation)
	var elevation uint32
	var returnLength uint32
	// syscall.TokenInformationClass = uint32; TokenElevation = 20
	err = syscall.GetTokenInformation(
		token,
		uint32(TOKEN_ELEVATION),
		(*byte)(unsafe.Pointer(&elevation)),
		uint32(unsafe.Sizeof(elevation)),
		&returnLength,
	)
	if err != nil {
		return false
	}

	// elevation != 0 表示管理员权限令牌已提升
	return elevation != 0
}

// RunAsAdmin 以管理员权限重新运行当前可执行文件
func RunAsAdmin() {
	exePath, err := os.Executable()
	if err != nil {
		return
	}

	// 拼接原参数
	args := strings.Join(os.Args[1:], " ")

	// 准备 Unicode 字符串指针
	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exePath)
	argsPtr, _ := syscall.UTF16PtrFromString(args)

	// 调用 ShellExecuteW(runas) 以管理员权限启动
	ret, _, _ := procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verbPtr)),
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(argsPtr)),
		0,
		SW_SHOWNORMAL,
	)

	// 大于 32 表示执行成功
	if ret > 32 {
		// 当前进程退出，让新的进程继续运行
		os.Exit(0)
	}
}
