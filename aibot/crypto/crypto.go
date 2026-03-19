// Package crypto 提供企业微信智能机器人 SDK 的加密解密功能
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// DecryptFile 解密文件数据
// 使用 AES-256-CBC 算法，数据采用 PKCS#7 填充
// aesKey 为 Base64 编码的密钥
func DecryptFile(encryptedData []byte, aesKey string) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("encrypted data cannot be empty")
	}

	if aesKey == "" {
		return nil, fmt.Errorf("aes key cannot be empty")
	}

	// 解码 Base64 密钥
	key, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 aes key: %w", err)
	}

	// 验证密钥长度
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid aes key length: expected 32 bytes, got %d bytes", len(key))
	}

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create aes cipher: %w", err)
	}

	// 验证加密数据长度
	if len(encryptedData) < aes.BlockSize {
		return nil, fmt.Errorf("encrypted data is too short: %d bytes, minimum is %d bytes",
			len(encryptedData), aes.BlockSize)
	}

	// 获取 IV（前 16 字节）
	iv := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]

	// 验证加密数据长度是块大小的倍数
	if len(encryptedData)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("encrypted data length is not a multiple of block size: %d bytes",
			len(encryptedData))
	}

	// 创建 CBC 解密器
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密数据（原地解密）
	decryptedData := make([]byte, len(encryptedData))
	mode.CryptBlocks(decryptedData, encryptedData)

	// 移除 PKCS#7 填充
	decryptedData, err = removePKCS7Padding(decryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to remove PKCS#7 padding: %w", err)
	}

	return decryptedData, nil
}

// removePKCS7Padding 移除 PKCS#7 填充
func removePKCS7Padding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data cannot be empty")
	}

	padding := int(data[len(data)-1])
	if padding == 0 || padding > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding size: %d", padding)
	}

	// 验证所有填充字节是否相同
	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return nil, fmt.Errorf("invalid padding byte at position %d: expected %d, got %d",
				i, padding, data[i])
		}
	}

	return data[:len(data)-padding], nil
}

// EncryptFile 加密文件数据（用于测试或本地加密）
// 使用 AES-256-CBC 算法，数据采用 PKCS#7 填充
// aesKey 为 Base64 编码的密钥
func EncryptFile(plainData []byte, aesKey string) ([]byte, error) {
	if len(plainData) == 0 {
		return nil, fmt.Errorf("plain data cannot be empty")
	}

	if aesKey == "" {
		return nil, fmt.Errorf("aes key cannot be empty")
	}

	// 解码 Base64 密钥
	key, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 aes key: %w", err)
	}

	// 验证密钥长度
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid aes key length: expected 32 bytes, got %d bytes", len(key))
	}

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create aes cipher: %w", err)
	}

	// 添加 PKCS#7 填充
	paddedData := addPKCS7Padding(plainData, aes.BlockSize)

	// 生成随机 IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate random iv: %w", err)
	}

	// 创建 CBC 加密器
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密数据
	encryptedData := make([]byte, len(paddedData))
	mode.CryptBlocks(encryptedData, paddedData)

	// 将 IV 添加到加密数据前
	result := make([]byte, len(iv)+len(encryptedData))
	copy(result[:len(iv)], iv)
	copy(result[len(iv):], encryptedData)

	return result, nil
}

// addPKCS7Padding 添加 PKCS#7 填充
func addPKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// GenerateAESKey 生成随机的 AES-256 密钥（Base64 编码）
func GenerateAESKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// ValidateAESKey 验证 AES 密钥格式
func ValidateAESKey(aesKey string) error {
	if aesKey == "" {
		return fmt.Errorf("aes key cannot be empty")
	}

	// 尝试解码 Base64
	key, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return fmt.Errorf("invalid base64 encoding: %w", err)
	}

	// 验证密钥长度
	if len(key) != 32 {
		return fmt.Errorf("invalid key length: expected 32 bytes, got %d bytes", len(key))
	}

	return nil
}

// DecryptFileWithIV 解密文件数据（指定 IV）
// 用于需要自定义 IV 的场景
func DecryptFileWithIV(encryptedData []byte, aesKey string, iv []byte) ([]byte, error) {
	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("encrypted data cannot be empty")
	}

	if aesKey == "" {
		return nil, fmt.Errorf("aes key cannot be empty")
	}

	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("invalid iv length: expected %d bytes, got %d bytes",
			aes.BlockSize, len(iv))
	}

	// 解码 Base64 密钥
	key, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 aes key: %w", err)
	}

	// 验证密钥长度
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid aes key length: expected 32 bytes, got %d bytes", len(key))
	}

	// 验证加密数据长度是块大小的倍数
	if len(encryptedData)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("encrypted data length is not a multiple of block size: %d bytes",
			len(encryptedData))
	}

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create aes cipher: %w", err)
	}

	// 创建 CBC 解密器
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密数据（原地解密）
	decryptedData := make([]byte, len(encryptedData))
	mode.CryptBlocks(decryptedData, encryptedData)

	// 移除 PKCS#7 填充
	decryptedData, err = removePKCS7Padding(decryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to remove PKCS#7 padding: %w", err)
	}

	return decryptedData, nil
}

// ExtractIVFromAESKey 从 AES 密钥中提取 IV（前 16 字节）
// 根据企业微信文档，IV 取 aeskey 前 16 字节
func ExtractIVFromAESKey(aesKey string) ([]byte, error) {
	if aesKey == "" {
		return nil, fmt.Errorf("aes key cannot be empty")
	}

	// 解码 Base64 密钥
	key, err := base64.StdEncoding.DecodeString(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 aes key: %w", err)
	}

	// 验证密钥长度
	if len(key) < aes.BlockSize {
		return nil, fmt.Errorf("aes key is too short: expected at least %d bytes, got %d bytes",
			aes.BlockSize, len(key))
	}

	// 提取前 16 字节作为 IV
	iv := make([]byte, aes.BlockSize)
	copy(iv, key[:aes.BlockSize])

	return iv, nil
}
