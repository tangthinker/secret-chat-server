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

### EncryptData

```c
char* EncryptData(const char* data);
```

- **参数**:
  - `data`: Base64 编码的输入数据
- **返回**: Base64 编码的加密数据，失败返回 NULL
- **注意**: 
  - 使用默认密钥进行加密
  - 返回的字符串需要使用 `FreeString` 释放

### DecryptData

```c
char* DecryptData(const char* data);
```

- **参数**:
  - `data`: Base64 编码的加密数据
- **返回**: Base64 编码的原始数据，失败返回 NULL
- **注意**: 
  - 使用默认密钥进行解密
  - 返回的字符串需要使用 `FreeString` 释放

### FreeString

```c
void FreeString(char* str);
```

- **参数**:
  - `str`: 需要释放的 C 字符串指针
- **功能**: 释放由 `EncryptData` 或 `DecryptData` 分配的内存

## Flutter 使用示例

```dart
import 'dart:ffi';
import 'package:ffi/ffi.dart';

// 加载库（macOS 示例）
final DynamicLibrary lib = DynamicLibrary.open('macos/Runner/Frameworks/libcrypto_arm64.dylib');

// 定义函数签名
typedef EncryptDataNative = Pointer<Utf8> Function(Pointer<Utf8> data);
typedef EncryptData = Pointer<Utf8> Function(Pointer<Utf8> data);

typedef DecryptDataNative = Pointer<Utf8> Function(Pointer<Utf8> data);
typedef DecryptData = Pointer<Utf8> Function(Pointer<Utf8> data);

typedef FreeStringNative = Void Function(Pointer<Utf8>);
typedef FreeString = void Function(Pointer<Utf8>);

// 获取函数
final encryptData = lib.lookupFunction<EncryptDataNative, EncryptData>(
    'EncryptData');
final decryptData = lib.lookupFunction<DecryptDataNative, DecryptData>(
    'DecryptData');
final freeString = lib.lookupFunction<FreeStringNative, FreeString>(
    'FreeString');

// 使用示例
String encrypt(String data) {
  final dataPtr = data.toNativeUtf8();
  
  try {
    final resultPtr = encryptData(dataPtr);
    if (resultPtr == nullptr) {
      throw Exception('Encryption failed');
    }
    final result = resultPtr.toDartString();
    freeString(resultPtr);
    return result;
  } finally {
    malloc.free(dataPtr);
  }
}

String decrypt(String data) {
  final dataPtr = data.toNativeUtf8();
  
  try {
    final resultPtr = decryptData(dataPtr);
    if (resultPtr == nullptr) {
      throw Exception('Decryption failed');
    }
    final result = resultPtr.toDartString();
    freeString(resultPtr);
    return result;
  } finally {
    malloc.free(dataPtr);
  }
}
```

## 加密算法

- **算法**: AES-256-GCM
- **密钥**: 使用默认密钥（`EncryptKey` 常量）
- **密钥派生**: SHA-256 哈希
- **数据格式**: Base64 编码
- **Nonce**: 随机生成，包含在密文中

## 注意事项

1. 所有输入和输出数据都使用 Base64 编码
2. 使用固定的默认密钥，通过 SHA-256 派生为 32 字节密钥
3. 必须在使用完返回的字符串后调用 `FreeString` 释放内存
4. 函数返回 NULL 表示操作失败
5. 密钥不可配置，统一使用默认密钥

