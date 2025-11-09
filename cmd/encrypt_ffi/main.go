package main

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"encoding/base64"
	"unsafe"

	"github.com/tangthinker/secret-chat-server/pkg"
)

// EncryptData 加密数据（C 导出函数）
// 参数：data - base64 编码的输入数据（使用默认密钥）
// 返回：base64 编码的加密数据，如果失败返回 NULL
//
//export EncryptData
func EncryptData(data *C.char) *C.char {
	if data == nil {
		return nil
	}

	// 将 C 字符串转换为 Go 字符串
	dataStr := C.GoString(data)

	// Base64 解码输入数据
	inputData, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return nil
	}

	// 加密数据（使用 pkg 包中的函数，使用默认密钥）
	encrypted, err := pkg.Encrypt(inputData)
	if err != nil {
		return nil
	}

	// Base64 编码输出数据
	outputStr := base64.StdEncoding.EncodeToString(encrypted)

	// 分配 C 字符串内存
	return C.CString(outputStr)
}

// DecryptData 解密数据（C 导出函数）
// 参数：data - base64 编码的加密数据（使用默认密钥）
// 返回：base64 编码的原始数据，如果失败返回 NULL
//
//export DecryptData
func DecryptData(data *C.char) *C.char {
	if data == nil {
		return nil
	}

	// 将 C 字符串转换为 Go 字符串
	dataStr := C.GoString(data)

	// Base64 解码输入数据
	inputData, err := base64.StdEncoding.DecodeString(dataStr)
	if err != nil {
		return nil
	}

	// 解密数据（使用 pkg 包中的函数，使用默认密钥）
	decrypted, err := pkg.Decrypt(inputData)
	if err != nil {
		return nil
	}

	// Base64 编码输出数据
	outputStr := base64.StdEncoding.EncodeToString(decrypted)

	// 分配 C 字符串内存
	return C.CString(outputStr)
}

// FreeString 释放 C 字符串内存（C 导出函数）
// 调用者必须在使用完字符串后调用此函数释放内存
//
//export FreeString
func FreeString(str *C.char) {
	if str != nil {
		C.free(unsafe.Pointer(str))
	}
}

// main 函数必须存在，但不会被调用（因为使用了 buildmode=c-shared）
func main() {}
