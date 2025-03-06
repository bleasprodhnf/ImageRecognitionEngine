package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptionService 提供数据加密和解密功能
type EncryptionService struct {
	key []byte
}

// NewEncryptionService 创建一个新的加密服务实例
func NewEncryptionService(key string) (*EncryptionService, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("加密密钥长度必须为16、24或32字节")
	}

	return &EncryptionService{
		key: []byte(key),
	}, nil
}

// Encrypt 使用AES-GCM加密数据
func (s *EncryptionService) Encrypt(plaintext string) (string, error) {
	// 创建cipher
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 创建nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// 返回Base64编码的密文
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密AES-GCM加密的数据
func (s *EncryptionService) Decrypt(ciphertext string) (string, error) {
	// 解码Base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 创建cipher
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 提取nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("密文长度不足")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]

	// 解密
	plaintextBytes, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintextBytes), nil
}

// EncryptSensitiveData 加密敏感数据
func (s *EncryptionService) EncryptSensitiveData(data map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	for key, value := range data {
		encrypted, err := s.Encrypt(value)
		if err != nil {
			return nil, err
		}
		result[key] = encrypted
	}

	return result, nil
}

// DecryptSensitiveData 解密敏感数据
func (s *EncryptionService) DecryptSensitiveData(data map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	for key, value := range data {
		decrypted, err := s.Decrypt(value)
		if err != nil {
			return nil, err
		}
		result[key] = decrypted
	}

	return result, nil
}