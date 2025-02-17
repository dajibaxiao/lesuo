package components

import (
	"golang.org/x/sys/windows/registry"
)

const (
	keyPath   = `Software\ZLesuo`
	valueName = "Machine"
)

// WriteRegistry 写入指定注册表键和值
func WriteRegistry(key registry.Key, path, name, value string) error {
	k, _, err := registry.CreateKey(key, path, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	err = k.SetStringValue(name, value)
	if err != nil {
		return err
	}
	return nil
}

// QueryRegistry 查询指定注册表键和值
func QueryRegistry(key registry.Key, path, name string) (string, error) {
	k, err := registry.OpenKey(key, path, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	value, _, err := k.GetStringValue(name)
	if err != nil {
		return "", err
	}
	return value, nil
}

// KeyExists 检查指定注册表键是否存在
func KeyExists(key registry.Key, path string) (bool, error) {
	k, err := registry.OpenKey(key, path, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, err
	}
	defer k.Close()
	return true, nil
}

// DeleteRegistry 删除指定注册表键
func DeleteRegistry(key registry.Key, path string) error {
	err := registry.DeleteKey(key, path)
	if err != nil {
		return err
	}
	return nil
}

func Chazhucebiao() bool {
	// 检查注册表键是否存在
	exists, err := KeyExists(registry.CURRENT_USER, keyPath)
	if err != nil {
		//fmt.Println("检查注册表键是否存在时出错:", err)
		return false
	}
	if exists {
		//fmt.Println("注册表键存在")
		return true
	} else {
		//fmt.Println("注册表键不存在")
		return false
	}
}

func Xiezhucebiao(valueData string) bool {
	// 检查注册表键是否存在
	err := WriteRegistry(registry.CURRENT_USER, keyPath, valueName, valueData)
	if err != nil {
		return false
	}
	return true
}

func Duzhucebiao() string {
	value, err := QueryRegistry(registry.CURRENT_USER, keyPath, valueName)
	if err != nil {
		return ""
	}
	return value
}

func CSdu() string {
	value, err := QueryRegistry(registry.CURRENT_USER, keyPath, "data")
	if err != nil {
		return ""
	}
	return value
}
func CSxie(valueData string) bool {
	// 检查注册表键是否存在
	err := WriteRegistry(registry.CURRENT_USER, keyPath, "data", valueData)
	if err != nil {
		return false
	}
	return true
}

func CSVIPdu() string {
	value, err := QueryRegistry(registry.CURRENT_USER, keyPath, "Key")
	if err != nil {
		return ""
	}
	return value
}
func CSVIPxie(valueData string) bool {
	// 检查注册表键是否存在
	err := WriteRegistry(registry.CURRENT_USER, keyPath, "Key", valueData)
	if err != nil {
		return false
	}
	return true
}
