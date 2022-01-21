package FamilyMartLogistics

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
)

func generateRequestBody(apiKey string, privateKey string, timestamp int64, xmlString string) string {
	sigDataStr := fmt.Sprint("ApiKey=",apiKey,"&Data=",xmlString,"&TimeStamp=",timestamp)
	str := fmt.Sprint("ApiKey=",apiKey,"&Data=",url.QueryEscape(xmlString),"&TimeStamp=",timestamp)
	h := hmac.New(sha512.New,[]byte(privateKey))
	h.Write([]byte(sigDataStr))
	sigData := h.Sum(nil)

	sig := hex.EncodeToString(sigData)

	//fmt.Println("EncodeData: ", sigDataStr);
	//fmt.Println("EncodeSign: ", sig);

	body := fmt.Sprint(str,"&Signature=",sig)
	return body
}