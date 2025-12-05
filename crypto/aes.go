package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// GenerateKey 生成一个 32 字节的安全随机密钥
// 对应 AES-256 标准
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	// crypto/rand 是真随机（CSPRNG），千万别用 math/rand
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// Encrypt 使用 AES-GCM 模式加密文本
// 返回：十六进制编码的密文
func Encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// GCM 模式：既加密又验证完整性，防止数据被篡改
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 创建随机数 Nonce (必须唯一，但在 GCM 中只需要随机即可)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Seal = 加密。
	// 第一个参数是 dst (目标切片)，如果我们传 nonce，它会把密文追加到 nonce 后面。
	// 这样做的好处是解密时方便切分：前 12 字节是 nonce，后面是密文。
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// 转成 Hex 字符串方便在网络传输
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt 解密
func Decrypt(encryptedHex string, key []byte) (string, error) {
	data, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// 拆分 Nonce 和 真正的密文
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Open = 解密
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// 如果密钥不对，或者数据被篡改，这里会直接报错
		return "", err
	}

	return string(plaintext), nil
}
