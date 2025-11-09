# Go 加密库 FFI 导出

这个目录包含用于编译为 C 共享库的 Go 加密代码，供 Flutter FFI 调用。

## 功能

- **EncryptData**: 使用 AES-256-GCM 加密数据
- **DecryptData**: 使用 AES-256-GCM 解密数据
- **FreeString**: 释放 C 字符串内存

## 编译

### 快速编译（当前平台）

```bash
cd cmd/encrypt_ffi
go build -buildmode=c-shared -o libcrypto.so main.go
```

### 编译所有平台

```bash
cd cmd/encrypt_ffi
chmod +x build.sh
./build.sh
```

编译后的库文件会输出到 `libs/` 目录。

## Flutter 项目集成（macOS）

### 1. 复制共享库到 Flutter 项目

将编译好的 `.dylib` 文件复制到 Flutter 项目的以下目录：

```
your_flutter_project/
  └── macos/
      └── Runner/
          └── Frameworks/
              └── libcrypto_arm64.dylib  (或 libcrypto_amd64.dylib)
```

**推荐做法：**
- 对于 Apple Silicon Mac (M1/M2/M3)：使用 `libcrypto_arm64.dylib`
- 对于 Intel Mac：使用 `libcrypto_amd64.dylib`
- 或者创建一个通用二进制文件（Universal Binary）包含两个架构

### 2. 在 Xcode 中配置

1. 打开 Flutter 项目的 Xcode workspace：
   ```bash
   open macos/Runner.xcworkspace
   ```

2. 在 Xcode 中配置：
   - 选择 `Runner` 项目 → `Runner` target
   - 进入 `Build Phases` 标签
   - 在 `Copy Bundle Resources` 中，点击 `+` 添加 `libcrypto_arm64.dylib`
   - 在 `Link Binary With Libraries` 中，确保库已添加（状态设为 `Optional`）

3. 设置库搜索路径：
   - 进入 `Build Settings` 标签
   - 搜索 `Library Search Paths`
   - 添加：`$(PROJECT_DIR)/Runner/Frameworks`

### 3. 在 Dart 代码中加载

```dart
import 'dart:ffi';
import 'dart:io';
import 'package:ffi/ffi.dart';
import 'package:path/path.dart' as path;

// 根据架构加载对应的库
DynamicLibrary loadCryptoLibrary() {
  if (Platform.isMacOS) {
    // 获取可执行文件所在目录
    final executablePath = Platform.resolvedExecutable;
    final executableDir = path.dirname(executablePath);
    
    // 尝试加载库（库应该在 Frameworks 目录中）
    final libraryPath = path.join(
      executableDir,
      '..',
      'Frameworks',
      'libcrypto_arm64.dylib'
    );
    
    return DynamicLibrary.open(libraryPath);
  }
  throw UnsupportedError('Platform not supported');
}

// 或者使用相对路径（开发时）
DynamicLibrary loadCryptoLibraryDev() {
  return DynamicLibrary.open('macos/Runner/Frameworks/libcrypto_arm64.dylib');
}
```

## C API

### BinaryData 结构体

```c
typedef struct {
    char* data;    // 数据指针
    size_t len;    // 数据长度
} BinaryData;
```

### EncryptData

```c
BinaryData EncryptData(const char* data, size_t dataLen);
```

- **参数**:
  - `data`: 输入数据的指针（二进制数据）
  - `dataLen`: 输入数据的长度
- **返回**: `BinaryData` 结构体，包含加密数据的指针和长度
  - 成功：`data` 指向加密数据，`len` 为数据长度
  - 失败：`data` 为 NULL，`len` 为 0
- **注意**: 
  - 使用默认密钥进行加密
  - 返回的数据需要使用 `FreeBinaryData` 释放

### DecryptData

```c
BinaryData DecryptData(const char* data, size_t dataLen);
```

- **参数**:
  - `data`: 加密数据的指针（二进制数据）
  - `dataLen`: 加密数据的长度
- **返回**: `BinaryData` 结构体，包含解密数据的指针和长度
  - 成功：`data` 指向解密数据，`len` 为数据长度
  - 失败：`data` 为 NULL，`len` 为 0
- **注意**: 
  - 使用默认密钥进行解密
  - 返回的数据需要使用 `FreeBinaryData` 释放

### FreeBinaryData

```c
void FreeBinaryData(BinaryData* binaryData);
```

- **参数**:
  - `binaryData`: 需要释放的 `BinaryData` 结构体指针
- **功能**: 释放由 `EncryptData` 或 `DecryptData` 分配的内存

## Flutter 使用示例

```dart
import 'dart:ffi';
import 'dart:typed_data';
import 'dart:convert';
import 'package:ffi/ffi.dart';

// 加载库（macOS 示例）
final DynamicLibrary lib = DynamicLibrary.open('macos/Runner/Frameworks/libcrypto_arm64.dylib');

// 定义 BinaryData 结构体
class BinaryData extends Struct {
  @IntPtr()
  external Pointer<Uint8> data;
  
  @Int64()
  external int len;
}

// 定义函数签名
typedef EncryptDataNative = BinaryData Function(Pointer<Uint8> data, Int64 dataLen);
typedef EncryptData = BinaryData Function(Pointer<Uint8> data, int dataLen);

typedef DecryptDataNative = BinaryData Function(Pointer<Uint8> data, Int64 dataLen);
typedef DecryptData = BinaryData Function(Pointer<Uint8> data, int dataLen);

typedef FreeBinaryDataNative = Void Function(Pointer<BinaryData>);
typedef FreeBinaryData = void Function(Pointer<BinaryData>);

// 获取函数
final encryptData = lib.lookupFunction<EncryptDataNative, EncryptData>(
    'EncryptData');
final decryptData = lib.lookupFunction<DecryptDataNative, DecryptData>(
    'DecryptData');
final freeBinaryData = lib.lookupFunction<FreeBinaryDataNative, FreeBinaryData>(
    'FreeBinaryData');

// 使用示例
Uint8List encrypt(Uint8List data) {
  final dataPtr = malloc<Uint8>(data.length);
  dataPtr.asTypedList(data.length).setAll(0, data);
  
  try {
    final result = encryptData(dataPtr, data.length);
    if (result.data == nullptr || result.len == 0) {
      throw Exception('Encryption failed');
    }
    
    // 复制结果数据
    final encrypted = result.data.asTypedList(result.len).toList();
    
    // 释放内存
    final resultPtr = malloc<BinaryData>()..ref = result;
    freeBinaryData(resultPtr);
    malloc.free(resultPtr);
    
    return Uint8List.fromList(encrypted);
  } finally {
    malloc.free(dataPtr);
  }
}

Uint8List decrypt(Uint8List data) {
  final dataPtr = malloc<Uint8>(data.length);
  dataPtr.asTypedList(data.length).setAll(0, data);
  
  try {
    final result = decryptData(dataPtr, data.length);
    if (result.data == nullptr || result.len == 0) {
      throw Exception('Decryption failed');
    }
    
    // 复制结果数据
    final decrypted = result.data.asTypedList(result.len).toList();
    
    // 释放内存
    final resultPtr = malloc<BinaryData>()..ref = result;
    freeBinaryData(resultPtr);
    malloc.free(resultPtr);
    
    return Uint8List.fromList(decrypted);
  } finally {
    malloc.free(dataPtr);
  }
}

// 如果需要处理字符串，可以添加辅助函数
String encryptString(String text) {
  final data = utf8.encode(text);
  final encrypted = encrypt(Uint8List.fromList(data));
  return base64Encode(encrypted); // 如果需要 base64 编码
}

String decryptString(String encryptedBase64) {
  final encrypted = base64Decode(encryptedBase64);
  final decrypted = decrypt(Uint8List.fromList(encrypted));
  return utf8.decode(decrypted);
}
```

## 加密算法

- **算法**: AES-256-GCM
- **密钥**: 使用默认密钥（`EncryptKey` 常量）
- **密钥派生**: SHA-256 哈希
- **数据格式**: 二进制数据（直接处理字节）
- **Nonce**: 随机生成，包含在密文中

## 注意事项

1. 所有输入和输出数据都是二进制数据（字节数组），不再使用 Base64 编码
2. 使用固定的默认密钥，通过 SHA-256 派生为 32 字节密钥
3. 必须在使用完返回的数据后调用 `FreeBinaryData` 释放内存
4. 函数返回 `data` 为 NULL 或 `len` 为 0 表示操作失败
5. 密钥不可配置，统一使用默认密钥
6. 输入数据长度通过 `dataLen` 参数传递，输出数据长度通过返回结构体的 `len` 字段获取

