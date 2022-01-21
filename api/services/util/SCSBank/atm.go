package SCSBank

import (
	"api/services/dao/sequence"
	"api/services/util/log"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GenerateAtmVirtualAccount(prefixCode string, money, afterDayCount int) (string, error) {
	today := time.Now()
	dateLitmit := today.YearDay() + afterDayCount
	year := (today.Year() % 10)
	str := fmt.Sprintf("%d%d",year,dateLitmit)
	seq, err := sequence.GetKgiVirtualSeq()
	if err != nil {
		return "", err
	}
	return createAtmVirtualAccount(prefixCode, seq,str,money), nil
}

func createAtmVirtualAccount(prefix string, seqStr string, dateStr string, money int) (string) {
	if len(seqStr) < 7 {
		seqStr = strings.Repeat("0",7 - len(seqStr)) + seqStr
	}
	preAccount := prefix + dateStr + seqStr
	account := numberStringToIntArray(preAccount)
	moneyValue := numberToIntArray(money)
	VirtualAccount := preAccount + checkSum(moneyValue,account)
	return VirtualAccount
}

func numberStringToIntArray(str string) (arr [12]int) {
	log.Info("numberStringToIntArray str", str)
	//03240000016
	l := len(str)
	account := str[l-12:l]

	for i, rune := range account {
		arr[i] ,_ = strconv.Atoi(string(rune))
	}

	return
}

func numberToIntArray(num int) (arr [12]int) {
	count := digitCountOfInt(num)
	tempV := num
	for i := 0; i < count; i++ {
		arr[11-i] = tempV % 10
		tempV /= 10
	}
	return
}

func digitCountOfInt(number int) (count int) {
	for number != 0 {
		number /= 10
		count += 1
	}
	return count
}

func checkSum(money, account [12]int) (sum string) {
	checkValue := [12]int{3,1,9,7,3,1,9,7,3,1,9,7}
	//checkValue := [12]int{1,9,8,7,6,5,4,3,2,1,9,8}
	tempValueC,tempValueF := 0,0

	for i := 0; i < 12; i++ {
		tempValueC += (account[i] * checkValue[i]) % 10
		tempValueF += (money[i] * checkValue[i]) % 10
	}

	finalCheckValue := (tempValueC+ tempValueF) % 10
	sum = strconv.Itoa(finalCheckValue)
	return
}