package entity

import (
	"api/services/Enum"
	"api/services/VO/OrderVo"
	"api/services/VO/Request"
	"api/services/util/tools"
	"time"
)

type MemberCardData struct {
	CardId       string    `xorm:"pk varchar(50) unique comment('序號')"`
	MemberId     string    `xorm:"varchar(50) notnull comment('會員編號')"`
	BankName     string    `xorm:"varchar(20) comment('銀行名稱')"`
	CardType     string    `xorm:"varchar(10) notnull comment('卡別')"`
	IsDebit      int64     `xorm:"tinyint(1) default 0 comment('是否為DEBIT')"`
	First4Digits string    `xorm:"varchar(4) notnull comment('前4碼')"`
	Last4Digits  string    `xorm:"varchar(4) notnull comment('後4碼')"`
	CardNumber   string    `xorm:"varchar(300) notnull comment('卡號 加密後')"`
	ExpiryDate   string    `xorm:"varchar(4) notnull comment('卡號到期日')"`
	Frequency    int64     `xorm:"int(10) default 0 comment('使用數')"`
	Status       string    `xorm:"varchar(20) default 'SUCCESS' notnull comment('狀態')"`
	DefaultCard  string    `xorm:"tinyint(1) default 0 comment('是否為預設')"`
	IsRelease    int64     `xorm:"tinyint(1) default 0 comment('是否請過款')"`
	IsForeign    int64     `xorm:"tinyint(1) default 0 comment('是否為國外')"`
	CreateTime   time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime   time.Time `xorm:"datetime notnull comment('更新時間')"`
}

func (c *MemberCardData) GenerateAuthRequest(vo OrderVo.CreditPaymentVo, params *Request.PayParams) AuthRequest {
	var AuthParams AuthRequest
	AuthParams.OrderId = vo.OrderId
	AuthParams.SellerId = vo.SellerId
	AuthParams.BuyerId = vo.BuyerId
	AuthParams.CardId = c.CardId
	number, _ := tools.AwsKMSDecrypt(c.CardNumber)
	AuthParams.CardNumber = number
	AuthParams.ExpireDate = c.ExpiryDate
	AuthParams.Security = params.CardSecurity
	AuthParams.Amount = int64(vo.TotalAmount)
	AuthParams.Type = vo.OrderType
	AuthParams.AuditStatus = vo.AuditStatus
	//是否使用次特店代碼
	AuthParams.IsMerchant = false
	//是否使用3D驗證 目前都設為3D驗證
	AuthParams.IsFirst = true
	//if c.IsRelease == 1 {
	//	AuthParams.IsFirst = false
	//}
	return AuthParams
}

type AuthRequest struct {
	OrderId     string
	SellerId    string
	BuyerId     string
	CardId      string
	CardNumber  string
	ExpireDate  string
	Security    string
	Amount      int64
	Type        string
	IsFirst     bool
	AuditStatus string
	MerchantId  string
	TerminalId  string
	IsMerchant  bool
}

//建立刷卡記錄
func (a *AuthRequest) GenerateGwCreditAuthData(ip, MerchantId, TerminalId string) GwCreditAuthData {
	var data GwCreditAuthData
	data.CardId = a.CardId
	data.MerchantId = MerchantId
	data.TerminalId = TerminalId
	data.OrderId = a.OrderId
	data.TramsAmount = a.Amount
	data.TransTime = time.Now()
	data.TransType = Enum.CreditTransTypeAuth
	data.TransStatus = Enum.CreditTransStatusInit
	data.PayType = a.Type
	data.CreditType = Enum.OrderTrans3D
	if a.IsFirst == true {
		data.AuditStatus = a.AuditStatus
	} else {
		data.AuditStatus = Enum.CreditAuditWait
	}
	if !a.IsFirst {
		data.CreditType = Enum.OrderTransN3D
	}
	data.CaptureStatus = Enum.CreditCaptureInit
	data.AuthIp = ip
	return data
}
