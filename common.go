package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type configurationFile struct {
	TelegramBotKey string
	Secret         string
	Host           string
	ListenAddress  string
	RelayTo        string
	UUID           string
	Jobs           []job
}

func readConfigurationFromFile(fileToRead string) (string, string, string, string, string, string, []job) {
	b, err := os.ReadFile(fileToRead)
	if err != nil {
		fmt.Print(err)
	}
	configContent := string(b)
	fmt.Println(configContent)
	var startupConfig configurationFile
	_ = json.Unmarshal([]byte(configContent), &startupConfig)
	return startupConfig.Host, startupConfig.Secret, startupConfig.TelegramBotKey, startupConfig.ListenAddress, startupConfig.RelayTo, startupConfig.UUID, startupConfig.Jobs
}

func GetAESDecrypted(encrypted string) ([]byte, error) {
	key := secretKey
	iv := secretIV

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)
	return ciphertext, nil
}

// PKCS5UnPadding  pads a certain blob of data with necessary data to be used in AES block cipher
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

// GetAESEncrypted encrypts given text in AES 256 CBC
func GetAESEncrypted(plaintext string) string {
	key := secretKey
	iv := secretIV

	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return ""
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)
	return str
}

func getKeyAndIV(secret string) (string, string) {
	h := sha256.New()
	h.Write([]byte(secret))
	sha256_hash := hex.EncodeToString(h.Sum(nil))
	return sha256_hash[:32], sha256_hash[32:48]
}

func sendData(data string) {
	if relayTo == "telegram" {
		sendDataToTelegram(data)
	} else if relayTo == "probe" {
		sendDataOverTCP(data)
	} else {
		sendDataOverTCP(data)
	}
}
