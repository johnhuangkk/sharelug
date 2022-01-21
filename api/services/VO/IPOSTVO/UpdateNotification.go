package IPOSTVO

import (
	"encoding/json"
)

type ShipStatusNotify struct {
	CheckMacValue string `json:"CheckMacValue"`
	info
}

type info struct {
	Chuge string `json:"Chuge"` // 儲格大小
	EcOrderNo string `json:"EcOrderNo"` // 電商訂單號碼
	EventTime string `json:"EventTime"` // 電商訂單號碼
	IBoxAddress string `json:"IBoxAddress"` // i 郵箱地址
	IBoxEmapUrl string `json:"IBoxEmapUrl"` // 電子地圖網址
	IBoxName string `json:"IBoxName"` // i 郵箱名稱
	MailNo string `json:"MailNo"` // 郵件編號
	PickupDeadline string `json:"PickupDeadline"` // 取件期限
	PickupPassword string `json:"PickupPassword"` // 取件密碼
	PickupStartTime string `json:"PickupStartTime"` // 招領起始時間
	Postage int64 `json:"Postage"` // 郵資
	PostOfficeName string `json:"PostOfficeName"` // 招領郵局名稱
	PostOfficeTel string `json:"PostOfficeTel"` // 招領郵局電話
	ReceiverPhonno string `json:"ReceiverPhonno"` // 收件人手機

	/**
	到達買家 i 郵箱(20)
	i 郵箱收寄(30)
	買家 i 郵箱逾期退賣家(70)
	買家 i 郵箱取件成功(80)
	到達賣家 i 郵箱(120)
	賣家 i 郵箱逾期退賣家 (170)
	賣家 i 郵箱取件成功(180)
	*/
	ShipStatus string `json:"ShipStatus"` // 郵件狀態代碼
	SysAdmId int64 `json:"SysAdmId"` // i 郵箱編號
	Timestamp string `json:"Timestamp"` // 招領起始時間
	VipId string `json:"VipId"` // 特約戶編號
	ZipCode string `json:"ZipCode"` // i 郵箱區碼
}

type VerifyNotify struct {
	info
}

func (s *ShipStatusNotify) VerifyParams() VerifyNotify {
	jsonData, _ := json.Marshal(s)
	a := VerifyNotify{}
	_ = json.Unmarshal(jsonData, &a)

	return a
}