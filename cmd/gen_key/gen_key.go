package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

// 生成ed25519公钥和私钥
// 打印输出
func main() {
	// 生成密钥对
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		fmt.Printf("生成密钥失败: %v\n", err)
		return
	}

	// 将公钥和私钥编码为 base64 格式
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey)
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKey)

	// 打印输出
	fmt.Println("=== Ed25519 密钥对 ===")
	fmt.Printf("公钥 (Public Key):\n%s\n\n", publicKeyBase64)
	fmt.Printf("私钥 (Private Key):\n%s\n", privateKeyBase64)
}
