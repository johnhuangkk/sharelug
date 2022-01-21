package SCSBank

import (
	"fmt"
	"strconv"
	"strings"
)

type IncomeRecord struct {
	RAccount string // 收款帳號 (index)
	SAccount string // 轉出帳號
	BTXDate  string // 索引日 (index)
	BUSDate  string // 營業日
	SeqNo    string // 序號 (index)
	TXDate   string // 交易日期yyyymmdd
	TXTime   string // 交易時間hhmmss
	VAccount string // 虛擬帳號
	Amount   int64  // 繳款金額
	DC       string // 借貸別
	Notes    string // 附註(交易類別)
	Valid    bool   // 資料是否有效
	Checksum string // 檢核碼
	Raw      string // 原始資料
}

func NewIncomeRecord(content string) (record IncomeRecord) {
	list := strings.Split(content,":")
	for i,_ := range list {
		list[i] = strings.TrimSpace(list[i])
	}

	record.Raw      = content
	record.RAccount = list[0]
	record.BTXDate  = list[1]
	record.SeqNo    = list[2]
	record.TXDate   = list[3]
	record.TXTime   = list[4]
	record.VAccount = list[5]

	amount,_ := strconv.ParseInt(list[6],10,0)
	record.Amount   = amount
	record.DC       = list[7]
	record.Notes    = list[8]
	record.SAccount = list[9]
	record.BUSDate  = list[10]
	record.Checksum = list[11]
	checksum,_ := strconv.Atoi(list[11])

	var count int32 = 0
	for _,v := range strings.TrimSuffix(content,list[11]) {
		count += v
	}
	fmt.Println("CheckSum:",list[11], " Value:",count)

	// 確認資料正確性
	record.Valid = (len(list) == 12) && (int(count) == checksum)
	fmt.Println(record)
	return
}