package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"math/rand"
)

const (
	randChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" //隨機字串字符源
)

var aesIV []byte  //固定16
var aesKey string //固定32

//產生固定大小的隨機字串
func randString(number int) string {
	cache := make([]byte, number)
	for index := range cache {
		cache[index] = randChars[rand.Intn(len(randChars))]
	}
	return string(cache)
}

//加密,注意：加密前解密key將會改變,故每次加密後必許等待解密,否則前面密文將失效永遠無法再解密
func Encode(inputText string) (outputCode string, err error) {

	//需要去加密的字串
	plaintext := []byte(inputText)

	//隨機aesKey
	aesKey = randString(32)

	//建立加密演算法 aes
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", fmt.Errorf("Error: NewCipher(%d bytes) = %s", len(aesKey), err)
	}

	//加密字串
	if len(aesIV) != 16 {
		aesIV = []byte(randString(16)) //隨機IV
	}

	cfb := cipher.NewCFBEncrypter(block, aesIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return fmt.Sprintf("%x", ciphertext), nil

}

//解密
func Decode(inputCode string) (outputText string, err error) {

	byteString, err := hex.DecodeString(inputCode)
	if err != nil {
		return "", fmt.Errorf("Error code : %v = %s", inputCode, err)
	}
	ciphertext := []byte(byteString)

	if len(aesIV) != 16 {
		aesIV = []byte("1234567890123456")
	}

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", fmt.Errorf("Error: NewCipher(%d bytes) = %s", len(aesKey), err)
	}

	// 解密字串
	cfbdec := cipher.NewCFBDecrypter(block, aesIV)
	plaintextCopy := make([]byte, len(ciphertext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)
	return string(plaintextCopy), nil
}
