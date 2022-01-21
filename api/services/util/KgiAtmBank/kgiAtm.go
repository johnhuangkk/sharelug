package KgiAtmBank

import (
	"api/services/Enum"
	"api/services/dao/sequence"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

func GenerateAtmVirtualAccount(money int, transType string) (string, error) {

	Money := tools.StringPadLeft(strconv.Itoa(money), 8)
	bankCode := viper.GetString("KgiBank.C2C.AtmCode")
	if transType == Enum.OrderTransB2c {
		bankCode = viper.GetString("KgiBank.B2C.AtmCode")
	}
	seq, err := sequence.GetKgiVirtualSeq()
	if err != nil {
		return "", err
	}
	no := tools.StringPadLeft(seq, 6)
	day := GetDayOfYear()
	log.Debug("", bankCode, no, day)
	account := fmt.Sprintf("%s%s%s", bankCode, no, day)
	code := GetCheckCode(account, Money)
	log.Debug("", account, code)
	VirtualAccount := fmt.Sprintf("%s%s", account, code)

	return VirtualAccount, nil
}

func GetAcctCheckCode(account string) int {
	sum := 0
	Account := strings.Split(account, "")
	Weight := strings.Split("371371371371371", "")

	for k, v := range Account {
		a, _ := strconv.Atoi(v)
		b, _ := strconv.Atoi(Weight[k])
		c := strconv.Itoa(a * b)
		d, _ := strconv.Atoi(c[len(c)-1:])
		sum += d
	}
	e := strconv.Itoa(sum)
	f, _ := strconv.Atoi(e[len(e)-1:])
	return f
}

func GetMoneyCheckCode(money string) int {
	sum := 0
	Money := strings.Split(money, "")
	Weight := strings.Split("87654321", "")

	for k, v := range Money {
		a, _ := strconv.Atoi(v)
		b, _ := strconv.Atoi(Weight[k])
		c := strconv.Itoa(a * b)
		d, _ := strconv.Atoi(c)
		sum += d
	}
	e := strconv.Itoa(sum)
	f, _ := strconv.Atoi(e[len(e)-1:])
	return f
}

func GetCheckCode(account string, money string) string {

	a := GetAcctCheckCode(account)
	b := GetMoneyCheckCode(money)
	c := strconv.Itoa(a + b)
	d, _ := strconv.Atoi(c[len(c)-1:])
	f := strconv.Itoa(10 - d)
	return f[len(f)-1:]
}

func GetDayOfYear() string {
	now := time.Now()
	year, _, _ := now.Date()
	y := strconv.Itoa(year)
	day := now.YearDay()
	exp := tools.StringPadLeft(strconv.Itoa(day + 1), 3)
	date := fmt.Sprintf("%s%s", y[len(y)-1:], exp)
	return date
}