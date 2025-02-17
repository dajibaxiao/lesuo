package defaults

import (
	_ "embed"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const (
	SPI_SETDESKWALLPAPER = 20
	SPIF_UPDATEINIFILE   = 1
	SPIF_SENDCHANGE      = 2
)

func setRegistryForStretch() error {
	// 打开注册表项 HKEY_CURRENT_USER\Control Panel\Desktop
	key, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	// 设置 WallpaperStyle = "2" 和 TileWallpaper = "0"
	if err := key.SetStringValue("WallpaperStyle", "2"); err != nil {
		return err
	}
	if err := key.SetStringValue("TileWallpaper", "0"); err != nil {
		return err
	}
	return nil
}

func SetWallpaper(path string) error {
	user32 := syscall.NewLazyDLL("user32.dll")
	procSystemParametersInfo := user32.NewProc("SystemParametersInfoW")

	pathUTF16, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	ret, _, err := procSystemParametersInfo.Call(
		uintptr(SPI_SETDESKWALLPAPER),
		0,
		uintptr(unsafe.Pointer(pathUTF16)),
		uintptr(SPIF_UPDATEINIFILE|SPIF_SENDCHANGE),
	)

	if ret == 0 {
		return err
	}
	return nil
}

func Bizhi(wallpaperData []byte) {
	// 修改注册表设置为拉伸显示
	if err := setRegistryForStretch(); err != nil {
		fmt.Println("修改注册表失败：", err)
		return
	}

	// 将嵌入的壁纸数据写入临时文件
	tempDir := os.TempDir()
	tempFile := tempDir + "\\embedded_wallpaper.jpg"
	//err := ioutil.WriteFile(tempFile, wallpaperData, 0644)
	err := os.WriteFile(tempFile, wallpaperData, 0644)
	if err != nil {
		fmt.Println("写入临时壁纸文件失败：", err)
		return
	}
	defer os.Remove(tempFile) // 可选：程序结束后删除临时文件

	// 设置桌面壁纸
	if err := SetWallpaper(tempFile); err != nil {
		fmt.Println("设置壁纸失败：", err)
	} else {
		fmt.Println("壁纸设置成功，并已拉伸铺满全屏")
	}
}
