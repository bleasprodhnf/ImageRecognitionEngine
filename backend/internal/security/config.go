package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"

	"github.com/pkg/errors"
)

// Config 安全配置
type Config struct {
	EncryptionKey []byte
}

// NewConfig 创建安全配置
func NewConfig() (*Config, error) {
	key := os.Getenv("CONFIG_ENCRYPTION_KEY")
	if key == "" {
		return nil, errors.New("CONFIG_ENCRYPTION_KEY environment variable is not set")
	}

	return &Config{
		EncryptionKey: []byte(key),
	}, nil
}

// Encrypt 加密数据
func (c *Config) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.EncryptionKey)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// Decrypt 解密数据
func (c *Config) Decrypt(ciphertext []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(c.EncryptionKey)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}