package tools

import (
	"fmt"
	"strings"
)

// 前留幾碼 後留幾碼
func doMask(input string, before int, after int) string {
	output := input
	if len(input) > 0 {
		mask := strings.Repeat("*", len([]rune(input)) - (before + after) )
		begin := []rune(input)[:before]
		end := []rune(input)[len([]rune(input)) - after:]
		output = fmt.Sprintf("%s%s%s", string(begin), mask, string(end))
	}
	return output
}

func MaskerEMail(mail string) string {
	str := strings.Split(mail, "@")
	return fmt.Sprintf("%s@%s", doMask(str[0], 3, 0 ), str[1])
}
/**
 * 電話遮碼
 * 091****123
 */
func MaskerPhone(mobile string) string {
	return doMask(mobile, 3, 3)
}

/**
 * 電話遮碼
 * *******123
 */
func MaskerPhoneLater(mobile string) string {
	return doMask(mobile, 0, 3)
}

//銀行帳號遮碼  **********123
func MaskerBankAccount(account string) string {
	return doMask(account, 0, 4)
}

/**
 * 姓名遮碼.
 * 遮碼規則如下:
 * 兩字中文名: 黃明 -> 黃*
 * 三字中文名: 黃小明 -> 黃*明
 * 四字中文名: 歐陽小明 -> 歐陽*明
 * 五字中文名: 歐陽豬太郎 -> 歐陽**郎
 */
func MaskerName(name string) string {
	length := len([]rune(name))
	var result string
	if length > 3 {
		result = doMask(name, 2, 1)
	} else if length < 3 {
		result = doMask(name, 1, 0)
	} else {
		result = doMask(name, 1, 1)
	}
	return result
}

/**
 * 地址遮碼
 * 新北市金山區金金金路110號 => 新北市金山區****11號
 */
func MaskerAddress(address string) string {
	address = strings.Replace(address, ",", "", -1)
	if len([]rune(address)) > 12 {
		return doMask(address, 9, 4)
	}
	return doMask(address, 9, 0)
}

func MaskerAddressLater(address string) string {
	return doMask(address, 0, 4)
}
//發票隱碼
func MaskerInvoice(number string) string {
	return doMask(number, 8, 0)
}

