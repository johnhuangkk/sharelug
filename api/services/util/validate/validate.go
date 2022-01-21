package validate

import (
	"api/services/util/tools"
	"regexp"
	"strings"
)

//mobile verify
func VerifyMobileFormat(phone string) bool {
	regular := "^09[0-9]{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}

func IsVerifyEnglish(str string) bool {
	regular := "[a-zA-Z]+"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(str)
}

//驗證信用卡卡號
func IsValidCreditNumber(number string) bool {
	card := strings.Split(tools.Trim(number), "")
	checksum := 0
	if len(card) != 16 {
		return false
	}
	for k, v := range card {
		if k == 15 {
			d := checksum % 10
			if 10 - d != tools.ParseInt(v) {
				return false
			}
		}
		if k % 2 == 0 {
			digit := tools.ParseInt(v) * 2
			if digit < 10 {
				checksum += digit
			} else {
				checksum += digit - 9
			}
		} else {
			digit := tools.ParseInt(v)
			checksum += digit
		}
	}
	if (checksum % 10) != 0 {
		return false
	}
	return true
}