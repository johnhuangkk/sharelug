package SLAppleMessage

import (
	"log"
	"net/http"
	"net/url"
)

func SendShareLugAppleMessage(target string, verifyCode string) bool {
	client := http.Client{}
	value := url.Values{}
	value.Set("phone",target)
	value.Set("verifyCode",verifyCode)

	_,err := client.PostForm("http://10.0.1.101:8066/imessage",value)
	if err != nil {
		log.Println("[SendShareLugAppleMessage] Err:",err)
		return false
	}
	return true
}
