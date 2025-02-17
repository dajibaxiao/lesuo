package components

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// 这是你给出的 RSA 公钥 (PEM 格式)。同样必须保留 BEGIN/END 。
var publicKeyStr = `-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBAMnLxiV6RynyjZGYwmRKfrtJeLQwY/zWDqUKsdXndJBRDXvoniWd39n7
F2oJ6gkYtHbP1Cpcmq6YC7Zl/6omfyDG5WcTe9gMpCiWXHOK0gFxgVlWS2DUASvG
RYFV/T1Kc/Zem/V6y2i8BARzuuK03NEbjJlfIzAWE7RSyg8tv8J5AgMBAAE=
-----END RSA PUBLIC KEY-----`

func Rsajiami(Entext string) string {
	// 2. 从 PEM 字符串中解析公钥
	pubKey, err := parseRSAPublicKey([]byte(publicKeyStr))
	if err != nil {
		//log.Fatalf("解析公钥失败: %v", err)
	}

	// 3. 使用公钥加密
	plaintext := []byte(Entext)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plaintext)
	if err != nil {
		//log.Fatalf("加密失败: %v", err)
	}
	// 为了方便查看或传输，这里可以将密文字节用 Base64 转成字符串
	ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)
	return string(ciphertextBase64)
}

// parseRSAPublicKey 将 RSA 公钥 PEM (PKCS#1) 解析为 *rsa.PublicKey
func parseRSAPublicKey(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("公钥 PEM 解码失败")
	}
	if block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("公钥 PEM Block 类型错误: %s", block.Type)
	}
	// 解析 DER 编码的公钥
	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %v", err)
	}
	return pub, nil
}
