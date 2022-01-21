package entity

import (
	"api/services/VO/IPOSTVO"
	"encoding/json"
	"strings"
	"time"
)

/**
	郵局大宗貨態 && 即時通知
*/
type PostShippingStatus struct {
	Id             int    `xorm:"pk int(32) autoincr"`
	MailNo         string `xorm:"index varchar(20) notnull comment('郵件號碼') "`
	HandleTime     string `xorm:"varchar(30) notnull comment('處理日期時間')"`
	ShippingStatus string `xorm:"varchar(20) notnull comment('郵件狀態')"`
	Branch         string `xorm:"varchar(30)  notnull comment('處理局')"`
	Chuge          string `xorm:"varchar(30)  comment('儲格大小')"`
	Postage        int64  `xorm:"int(5) default 0  comment('郵資')"`
	StatusCode     string `xorm:"varchar(5)  comment('狀態碼')"`
	CreateTime     string `xorm:"datetime notnull comment('建立時間')"`
	Detail         string `xorm:"text  comment('詳細資訊')"`
}

func (p *PostShippingStatus) SetData(params IPOSTVO.ShipStatusNotify) {
	var mappingStatus = map[string]string{
		// 新增狀態 有可能需要增加判斷條件
		"G4":  "交投",
		"H451E":  "無收件人 i 郵箱退件",
		"Z8":  "交投寄件人",
		"G9":  "寄件人 i 郵箱逾期轉招領",
		"I9":  "寄件人招領取件成功",

		"H451": "未妥投",
		"Z4":  "運輸途中",
		"Z1":  "信箱郵件轉運中",
		"H4":  "到宅投遞不成功",
		"G2":  "箱到宅收件人投遞不成功轉招領",
		"T2":  "招領逾期退寄件人",

		"20":  "到達買家 i 郵箱",
		"30":  "i 郵箱收寄",
		"70":  "買家 i 郵箱逾期退賣家",
		"75":  "買家 i 郵箱逾期取出逕退",
		"80":  "買家 i 郵箱取件成功",
		"120": "到達賣家 i 郵箱",
		"170": "賣家 i 郵箱逾期退賣家",
		"180": "賣家 i 郵箱取件成功",
	}

	p.MailNo = params.MailNo
	t, _ := time.Parse("2006-01-02 15:04:05", strings.Split(params.EventTime, ".")[0])
	p.HandleTime = t.Format(`2006-01-02 15:04:05`)
	p.ShippingStatus = mappingStatus[params.ShipStatus]
	p.CreateTime = params.Timestamp
	p.Chuge = params.Chuge
	p.StatusCode = params.ShipStatus
	p.Postage = params.Postage
	jsonData, _ := json.Marshal(params)
	p.Detail = string(jsonData)
}
