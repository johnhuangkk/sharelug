package Invoice

import (
	"api/services/VO/InvoiceVo"
	"api/services/util/tools"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"
)

func QRCodeProduct(vo InvoiceVo.QRCodeInvVo) string {
	var str []string
	for _, v := range vo.ProductArray {
		str = append(str, fmt.Sprintf("%s:%d:%d", v.ProductName, v.Quantity, v.ProductPrice))
	}
	return fmt.Sprintf("**%s", strings.Join(str, ":"))
}

func QRCodeINV(vo InvoiceVo.QRCodeInvVo) (string, error) {
	str := ""
	for _, v := range vo.ProductArray {
		str = str + fmt.Sprintf(":%s:%d:%d", v.ProductName, v.Quantity, v.ProductPrice)
	}
	plainText := fmt.Sprintf("%s%s0", vo.InvoiceNumber, vo.RandomNumber)
	cipherText := aesEncrypt(plainText, vo.AESKey)
	test := fmt.Sprintf("%s%s%s%s%s%s%s%s:**********:1:1:1%s",
		vo.InvoiceNumber,
		vo.InvoiceDate,
		vo.RandomNumber,
		vo.SalesAmount,
		//vo.TaxAmount,
		vo.TotalAmount,
		vo.BuyerIdentifier,
		//vo.RepresentIdentifier,
		vo.SellerIdentifier,
		//vo.BusinessIdentifier,
		cipherText, str)
	return test, nil
}

func aesEncrypt(plainText, Key string) string {
	//轉成byte
	OrigData := []byte(plainText)
	KeyData := tools.Hex2Bin(Key)
	//分组密鑰
	block, _ := aes.NewCipher(KeyData)
	//獲取密鑰的長度
	blockSize := block.BlockSize()
	//補全碼
	OrigData = tools.PKCS7Padding(OrigData, blockSize)
	//加密模式
	iv := tools.Hex2Bin("0EDF25C93A28D7B5FF5E45DA42F8A1B8")
	//iv,  _ = hex.DecodeString(string(iv))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//創建byte
	encrypted := make([]byte, len(OrigData))
	//加密
	blockMode.CryptBlocks(encrypted, OrigData)
	return base64.StdEncoding.EncodeToString(encrypted)
}