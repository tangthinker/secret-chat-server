#!/bin/bash

# 构建脚本：将 Go 代码编译为 C 共享库供 Flutter FFI 调用

# 不使用 set -e，因为某些平台的交叉编译可能会失败，这是正常的

OUTPUT_DIR="../../libs"
mkdir -p "$OUTPUT_DIR"

echo "Building for different platforms..."
echo "Note: Cross-compilation requires CGO and target platform toolchains"
echo ""

# 检查当前平台
CURRENT_OS=$(go env GOOS)
CURRENT_ARCH=$(go env GOARCH)

# macOS (arm64)
echo "Building for macOS (arm64)..."
if GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o "$OUTPUT_DIR/libcrypto_arm64.dylib" main.go 2>/dev/null; then
    echo "  ✓ macOS (arm64) build successful"
else
    echo "  ✗ macOS (arm64) build failed (may need macOS SDK)"
fi

# macOS (amd64)
echo "Building for macOS (amd64)..."
if GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o "$OUTPUT_DIR/libcrypto_amd64.dylib" main.go 2>/dev/null; then
    echo "  ✓ macOS (amd64) build successful"
else
    echo "  ✗ macOS (amd64) build failed (may need macOS SDK)"
fi

# Linux (amd64)
echo "Building for Linux (amd64)..."
if GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o "$OUTPUT_DIR/libcrypto_linux.so" main.go 2>/dev/null; then
    echo "  ✓ Linux (amd64) build successful"
else
    echo "  ✗ Linux (amd64) build failed (may need gcc cross-compiler)"
fi

# Windows (amd64)
echo "Building for Windows (amd64)..."
if GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o "$OUTPUT_DIR/libcrypto.dll" main.go 2>/dev/null; then
    echo "  ✓ Windows (amd64) build successful"
else
    echo "  ✗ Windows (amd64) build failed (may need MinGW cross-compiler)"
fi

# Android (arm64) - 需要 Android NDK
echo "Building for Android (arm64)..."
echo "  Note: Requires Android NDK and CC environment variable"
if [ -z "$CC" ] && [ -z "$ANDROID_NDK_HOME" ]; then
    echo "  ⚠ Skipping Android (arm64) - Android NDK not configured"
    echo "    Set ANDROID_NDK_HOME or CC to enable Android builds"
else
    if GOOS=android GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-shared -o "$OUTPUT_DIR/libcrypto_android_arm64.so" main.go 2>/dev/null; then
        echo "  ✓ Android (arm64) build successful"
    else
        echo "  ✗ Android (arm64) build failed"
    fi
fi

# Android (armv7) - 需要 Android NDK
echo "Building for Android (armv7)..."
if [ -z "$CC" ] && [ -z "$ANDROID_NDK_HOME" ]; then
    echo "  ⚠ Skipping Android (armv7) - Android NDK not configured"
else
    if GOOS=android GOARCH=arm CGO_ENABLED=1 go build -buildmode=c-shared -o "$OUTPUT_DIR/libcrypto_android_armv7.so" main.go 2>/dev/null; then
        echo "  ✓ Android (armv7) build successful"
    else
        echo "  ✗ Android (armv7) build failed"
    fi
fi

# iOS (arm64) - 需要 Xcode 和 iOS SDK
echo "Building for iOS (arm64)..."
echo "  Note: Requires Xcode and iOS SDK"
if [ "$CURRENT_OS" != "darwin" ]; then
    echo "  ⚠ Skipping iOS - requires macOS with Xcode"
elif ! xcodebuild -version &>/dev/null; then
    echo "  ⚠ Skipping iOS - Xcode not found"
else
    # iOS 需要使用 c-archive 而不是 c-shared
    if GOOS=ios GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=c-archive -o "$OUTPUT_DIR/libcrypto_ios_arm64" main.go 2>/dev/null; then
        echo "  ✓ iOS (arm64) build successful"
    else
        echo "  ✗ iOS (arm64) build failed (may need iOS SDK)"
    fi
fi

echo ""
echo "Build complete! Libraries are in $OUTPUT_DIR"
echo ""
echo "Generated files:"
ls -lh "$OUTPUT_DIR"/*.so "$OUTPUT_DIR"/*.dylib "$OUTPUT_DIR"/*.dll "$OUTPUT_DIR"/*.a 2>/dev/null || echo "No files generated"
echo ""
echo "Note: Some platforms may require additional setup:"
echo "  - Linux: Install gcc cross-compiler for cross-compilation"
echo "  - Windows: Install MinGW-w64 for cross-compilation"
echo "  - Android: Install Android NDK and set ANDROID_NDK_HOME"
echo "  - iOS: Requires macOS with Xcode installed"

