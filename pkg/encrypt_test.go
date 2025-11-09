package pkg

import (
	"bytes"
	"testing"
)

func TestEncrypt(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "正常加密 - 短文本",
			data:    []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "正常加密 - 长文本",
			data:    []byte("This is a longer text that contains more characters to test encryption with longer data."),
			wantErr: false,
		},
		{
			name:    "正常加密 - 二进制数据",
			data:    []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD},
			wantErr: false,
		},
		{
			name:    "正常加密 - JSON 数据",
			data:    []byte(`{"name":"test","value":123}`),
			wantErr: false,
		},
		{
			name:    "错误 - 空数据",
			data:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := Encrypt(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if encrypted == nil {
					t.Error("Encrypt() returned nil encrypted data")
					return
				}
				if len(encrypted) == 0 {
					t.Error("Encrypt() returned empty encrypted data")
					return
				}
				// 加密后的数据应该比原始数据长（因为包含 nonce）
				if len(encrypted) <= len(tt.data) {
					t.Errorf("Encrypt() encrypted data length (%d) should be greater than original (%d)", len(encrypted), len(tt.data))
				}
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "错误 - 空数据",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "错误 - 数据太短",
			data:    []byte{0x01, 0x02, 0x03},
			wantErr: true,
		},
		{
			name:    "错误 - 无效的密文",
			data:    make([]byte, 50), // 全零，不是有效的密文
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decrypted, err := Decrypt(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && decrypted == nil {
				t.Error("Decrypt() returned nil decrypted data")
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "短文本",
			data: []byte("Hello, World!"),
		},
		{
			name: "长文本",
			data: []byte("This is a longer text that contains more characters to test encryption with longer data."),
		},
		{
			name: "空字节",
			data: []byte{0x00, 0x00, 0x00},
		},
		{
			name: "二进制数据",
			data: []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD, 0xFC},
		},
		{
			name: "JSON 数据",
			data: []byte(`{"name":"test","value":123,"nested":{"key":"value"}}`),
		},
		{
			name: "Unicode 文本",
			data: []byte("你好，世界！🌍"),
		},
		{
			name: "单字节",
			data: []byte{0x42},
		},
		{
			name: "大块数据",
			data: bytes.Repeat([]byte("A"), 1024*10), // 10KB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 加密
			encrypted, err := Encrypt(tt.data)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}
			if encrypted == nil || len(encrypted) == 0 {
				t.Fatal("Encrypt() returned nil or empty encrypted data")
			}

			// 解密
			decrypted, err := Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}
			if decrypted == nil {
				t.Fatal("Decrypt() returned nil decrypted data")
			}

			// 验证解密后的数据与原始数据相同
			if !bytes.Equal(decrypted, tt.data) {
				t.Errorf("Decrypt(Encrypt(data)) = %v, want %v", decrypted, tt.data)
			}
		})
	}
}

func TestEncryptRandomness(t *testing.T) {
	// 测试多次加密相同数据，结果应该不同（因为 nonce 是随机的）
	data := []byte("test data")
	encrypted1, err := Encrypt(data)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	encrypted2, err := Encrypt(data)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// 两次加密的结果应该不同（因为 nonce 是随机的）
	if bytes.Equal(encrypted1, encrypted2) {
		t.Error("Encrypt() should produce different results for the same input (due to random nonce)")
	}

	// 但是解密后应该得到相同的结果
	decrypted1, err := Decrypt(encrypted1)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	decrypted2, err := Decrypt(encrypted2)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if !bytes.Equal(decrypted1, decrypted2) {
		t.Error("Decrypt() should produce the same result for different encryptions of the same data")
	}

	if !bytes.Equal(decrypted1, data) {
		t.Error("Decrypt() should produce the original data")
	}
}

func TestDecryptWrongData(t *testing.T) {
	// 测试用错误的密文解密
	originalData := []byte("test data")
	encrypted, err := Encrypt(originalData)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// 修改密文的一个字节
	corrupted := make([]byte, len(encrypted))
	copy(corrupted, encrypted)
	corrupted[0] ^= 0xFF // 翻转第一个字节

	// 解密应该失败
	_, err = Decrypt(corrupted)
	if err == nil {
		t.Error("Decrypt() should fail with corrupted ciphertext")
	}
}

func BenchmarkEncrypt(b *testing.B) {
	data := []byte("This is a test data for benchmarking encryption performance")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Encrypt(data)
		if err != nil {
			b.Fatalf("Encrypt() error = %v", err)
		}
	}
}

func BenchmarkDecrypt(b *testing.B) {
	data := []byte("This is a test data for benchmarking decryption performance")
	encrypted, err := Encrypt(data)
	if err != nil {
		b.Fatalf("Encrypt() error = %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decrypt(encrypted)
		if err != nil {
			b.Fatalf("Decrypt() error = %v", err)
		}
	}
}

func BenchmarkEncryptDecrypt(b *testing.B) {
	data := []byte("This is a test data for benchmarking encryption and decryption performance")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypted, err := Encrypt(data)
		if err != nil {
			b.Fatalf("Encrypt() error = %v", err)
		}
		_, err = Decrypt(encrypted)
		if err != nil {
			b.Fatalf("Decrypt() error = %v", err)
		}
	}
}
