package tools

import (
	"api/services/util/log"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/spf13/viper"
	"io"
	"strings"
)

//加密
func AesEncrypt(plainText, Key string) string {
	//轉成byte
	OrigData := []byte(plainText)
	KeyData := []byte(Key)
	//分组密鑰
	block, _ := aes.NewCipher(KeyData)
	//獲取密鑰的長度
	blockSize := block.BlockSize()
	//補全碼
	OrigData = PKCS7Padding(OrigData, blockSize)
	//加密模式
	blockMode := cipher.NewCBCEncrypter(block, KeyData[:blockSize])
	//創建byte
	encrypted := make([]byte, len(OrigData))
	//加密
	blockMode.CryptBlocks(encrypted, OrigData)
	return base64.StdEncoding.EncodeToString(encrypted)
}

//解密
func AesDecrypt(cipherText, Key string) string {
	//轉成byte
	EncryptedByte, _ := base64.StdEncoding.DecodeString(cipherText)
	KeyData := []byte(Key)
	//分组密鑰
	block, _ := aes.NewCipher(KeyData)
	//獲取密鑰的長度
	blockSize := block.BlockSize()
	//加密模式
	blockMode := cipher.NewCBCDecrypter(block, KeyData[:blockSize])
	// 创建数组
	plainText := make([]byte, len(EncryptedByte))
	// 解密
	blockMode.CryptBlocks(plainText, EncryptedByte)
	// 去码
	plainText = PKCS7UnPadding(plainText)
	return string(plainText)
}

//補碼
func PKCS7Padding(CipherText []byte, BlockSize int) []byte {
	padding := BlockSize - len(CipherText) % BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(CipherText, padtext...)
}

//去碼
func PKCS7UnPadding(OrigData []byte) []byte {
	length := len(OrigData)
	unpadding := int(OrigData[length-1])
	return OrigData[:(length - unpadding)]
}

func Base64Encode(message []byte) []byte {
	b := make([]byte, base64.StdEncoding.EncodedLen(len(message)))
	base64.StdEncoding.Encode(b, message)
	return b
}

func Hex2Bin(s string) []byte {
	ret, _ := hex.DecodeString(s)
	return ret
}

// md5
func MD5(str string) string {
	md5Module := md5.New()
	_, _ = io.WriteString(md5Module, str)
	return fmt.Sprintf("%x", md5Module.Sum(nil))
}

// SHA512Mac
func SHA512Mac(str string, key string) string {
	mac := hmac.New(sha512.New, []byte(key))
	_,_ = io.WriteString(mac, str)
	return fmt.Sprintf("%x", mac.Sum(nil))
}

// SHA384Mac
func SHA384Mac(str string, key string) string {
	mac := hmac.New(sha512.New384, []byte(key))
	_,_ = io.WriteString(mac, str)
	return fmt.Sprintf("%x", mac.Sum(nil))
}

// SHA256
func SHA256(str string) string {
	h := sha256.New()
	h.Write([]byte(strings.ReplaceAll(str, "\\\\", "\\")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//HMAC-SHA256
func SHA256ByKey(secret, str string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Base64Encode
func Base64EncodeByString(plaintext string) string {
	return base64.StdEncoding.EncodeToString([]byte(plaintext))
}

// Base64Decode
func Base64DecodeByString(encodeString string) string {
	str, _ := base64.StdEncoding.DecodeString(encodeString)
	return string(str)
}


// aws kms session 建立
func awsSessKms() (*kms.KMS, error) {
	var svc *kms.KMS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	if err != nil {
		fmt.Println("Error", err)
		return svc, err
	}

	svc = kms.New(sess)

	return svc, nil
}

// 使用 kms 加密
func AwsKMSEncrypt(plaintext string) (string, error) {
	var CiphertextBlob string
	svc, err := awsSessKms()
	if err != nil {
		log.Error("Aws Session Error", err)
		return CiphertextBlob, nil
	}

	result, err := svc.Encrypt(&kms.EncryptInput{
		KeyId: aws.String(viper.GetString("AWS.KMS_KEY")),
		Plaintext: []byte(plaintext),
	})

	if err != nil {
		log.Error("AwsKMSEncrypt Error:", err)
		return CiphertextBlob, fmt.Errorf("%s [%s]", "AwsKMSEncrypt error encrypting data", err.Error())
	}
	CiphertextBlob = string(result.CiphertextBlob)
	return Base64EncodeByString(CiphertextBlob), nil
}

// 使用 kms 解密
func AwsKMSDecrypt(CiphertextBlob string) (string, error) {
	var plaintext string
	svc, err := awsSessKms()
	if err != nil {
		log.Error("AwsKMSDecrypt Error:", err)
		return CiphertextBlob, nil
	}
	result, err := svc.Decrypt(&kms.DecryptInput{KeyId: aws.String(viper.GetString("AWS.KMS_KEY")), CiphertextBlob: []byte(Base64DecodeByString(CiphertextBlob))})
	if err != nil {
		log.Error("AwsKMSDecrypt Error", err)
		return "", nil
	}
	plaintext = string(result.Plaintext)
	return plaintext, nil
}