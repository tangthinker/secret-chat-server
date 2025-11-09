package main

/*
#include <stdlib.h>
#include <string.h>

// 返回结构体，包含数据指针和长度
typedef struct {
    char* data;
    size_t len;
} BinaryData;
*/
import "C"
import (
	"unsafe"

	"github.com/tangthinker/secret-chat-server/pkg"
)

// EncryptData 加密数据（C 导出函数）
// 参数：data - 输入数据的指针，dataLen - 输入数据长度（使用默认密钥）
// 返回：包含加密数据指针和长度的结构体，如果失败返回 NULL 指针和 0 长度
//
//export EncryptData
func EncryptData(data *C.char, dataLen C.size_t) C.BinaryData {
	var result C.BinaryData
	result.data = nil
	result.len = 0

	if data == nil || dataLen == 0 {
		return result
	}

	// 将 C 字节数组转换为 Go 字节切片
	inputData := C.GoBytes(unsafe.Pointer(data), C.int(dataLen))

	// 加密数据（使用 pkg 包中的函数，使用默认密钥）
	encrypted, err := pkg.Encrypt(inputData)
	if err != nil {
		return result
	}

	// 分配内存并复制加密后的数据
	if len(encrypted) > 0 {
		result.data = (*C.char)(C.malloc(C.size_t(len(encrypted))))
		if result.data == nil {
			return result
		}
		C.memcpy(unsafe.Pointer(result.data), unsafe.Pointer(&encrypted[0]), C.size_t(len(encrypted)))
		result.len = C.size_t(len(encrypted))
	}

	return result
}

// DecryptData 解密数据（C 导出函数）
// 参数：data - 加密数据的指针，dataLen - 加密数据长度（使用默认密钥）
// 返回：包含解密数据指针和长度的结构体，如果失败返回 NULL 指针和 0 长度
//
//export DecryptData
func DecryptData(data *C.char, dataLen C.size_t) C.BinaryData {
	var result C.BinaryData
	result.data = nil
	result.len = 0

	if data == nil || dataLen == 0 {
		return result
	}

	// 将 C 字节数组转换为 Go 字节切片
	inputData := C.GoBytes(unsafe.Pointer(data), C.int(dataLen))

	// 解密数据（使用 pkg 包中的函数，使用默认密钥）
	decrypted, err := pkg.Decrypt(inputData)
	if err != nil {
		return result
	}

	// 分配内存并复制解密后的数据
	if len(decrypted) > 0 {
		result.data = (*C.char)(C.malloc(C.size_t(len(decrypted))))
		if result.data == nil {
			return result
		}
		C.memcpy(unsafe.Pointer(result.data), unsafe.Pointer(&decrypted[0]), C.size_t(len(decrypted)))
		result.len = C.size_t(len(decrypted))
	}

	return result
}

// FreeBinaryData 释放二进制数据内存（C 导出函数）
// 调用者必须在使用完数据后调用此函数释放内存
//
//export FreeBinaryData
func FreeBinaryData(binaryData *C.BinaryData) {
	if binaryData != nil && binaryData.data != nil {
		C.free(unsafe.Pointer(binaryData.data))
		binaryData.data = nil
		binaryData.len = 0
	}
}

// main 函数必须存在，但不会被调用（因为使用了 buildmode=c-shared）
func main() {}
