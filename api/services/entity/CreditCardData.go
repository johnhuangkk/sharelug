package entity

import (
	"api/services/Enum"
	"api/services/VO/Request"
	"time"
)

type GwCreditAuthData struct {
	AuthId        int       `xorm:"pk int(10) unique autoincr comment('信用卡文易序號')"`
	MerchantId    string    `xorm:"varchar(15) notnull default '1300816653' comment('商店代碼')"`
	TerminalId    string    `xorm:"varchar(15) notnull default 'A2300432' comment('端末機代號')"`
	OrderId       string    `xorm:"varchar(100) notnull comment('訂單編號')"`
	TransType     string    `xorm:"varchar(10) notnull comment('交易類型')"`
	TramsAmount   int64     `xorm:"int(10) notnull comment('交易金額')"`
	CardId        string    `xorm:"varchar(50) notnull comment('交易卡片編號')"`
	PayType       string    `xorm:"varchar(20) notnull comment('交易類型')"`
	CreditType    string    `xorm:"varchar(10) notnull default '3D' comment('刷卡方式')"`
	TransTime     time.Time `xorm:"datetime notnull comment('交易時間')"`
	AuditStatus   string    `xorm:"varchar(10) notnull default 'INIT' comment('審核狀態')"`
	AuditStaff    string    `xorm:"varchar(10) comment('處理人員')"`
	AuditTime     time.Time `xorm:"datetime comment('審單時間')"`
	TransStatus   string    `xorm:"varchar(10) notnull comment('交易狀態')"`
	ApproveCode   string    `xorm:"varchar(10) comment('授權碼')"`
	ResponseCode  string    `xorm:"varchar(10) comment('回應代碼')"`
	ResponseMsg   string    `xorm:"text comment('回應訊息')"`
	BatchId       int       `xorm:"int(10) comment('批次ID')"`
	BatchTime     time.Time `xorm:"datetime comment('批次時間')"`
	CaptureStatus string    `xorm:"varchar(10) default 'INIT' comment('請款狀態')"`
	CaptureCode   string    `xorm:"varchar(10) comment('請款回應代碼')"`
	CaptureMsg    string    `xorm:"varchar(50) comment('請款回應訊息')"`
	Memo          string    `xorm:"text comment('審單備註')"`
	AuthIp        string    `xorm:"varchar(30) comment('交易IP')"`
	NoteTime      time.Time `xorm:"datetime comment('照會時間')"`
	PendingTime   time.Time `xorm:"datetime comment('待決時間')"`
	RefusedTime   time.Time `xorm:"datetime comment('拒絕時間')"`
	ReleaseTime   time.Time `xorm:"datetime comment('放行時間')"`
	CaptureTime   time.Time `xorm:"datetime comment('請款回應時間')"`
	CreateTime    time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime    time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type GwCreditAuthLog struct {
	Id         int       `xorm:"pk int(10) unique autoincr comment('序號')"`
	Response   string    `xorm:"text comment('回應訊息')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}

type CreditBatchRequestData struct {
	BatchId     int       `xorm:"pk int(6) unique autoincr comment('批次id')"`
	SendDate    string    `xorm:"varchar(10) notnull comment('送檔日期')"`
	Count       string    `xorm:"varchar(10) notnull comment('總筆數')"`
	Symbol      string    `xorm:"varchar(2) notnull comment('金額正負號')"`
	ToTalAmount string    `xorm:"varchar(10) notnull comment('總金額')"`
	Status      string    `xorm:"varchar(10) notnull comment('狀態')"`
	CreateTime  time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime  time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type QueryResponse struct {
	TransDate       string `json:"TransDate"`       //交易日期
	TransTime       string `json:"TransTime"`       //交易時間
	ApproveCode     string `json:"ApproveCode"`     //授權碼
	ResponseCode    string `json:"ResponseCode"`    //授權回應碼
	ResponseMsg     string `json:"ResponseMsg"`     //授權回應訊息
	InstallmentType string `json:"InstallmentType"` //分期手續費計價方式
	FirstAmt        string `json:"FirstAmt"`        //首期金額
	EachAmt         string `json:"EachAmt"`         //每期金額
	Fee             string `json:"Fee"`             //分期手續費
	RedeemType      string `json:"RedeemType"`      //紅利折抵方式
	RedeemUsed      string `json:"RedeemUsed"`      //紅利折抵點數
	RedeemBalance   string `json:"RedeemBalance"`   //紅利餘額
	CreditAmt       string `json:"CreditAmt"`       //持卡人自付金額
	RiskMark        string `json:"RiskMark"`        //風險卡號註記
	Foreign         string `json:"Foreign"`         //國外卡
	SecureStatus    string `json:"SecureStatus"`    //3D 認證結果
	OrderId         string `json:"OrderId"`
	Installment     string `json:"Installment"`
	MerchantID      string `json:"MerchantID"`
	PrivateData     string `json:"PrivateData"`
	TransMode       string `json:"TransMode"`
	TransAmt        string `json:"TransAmt"`
	TerminalID      string `json:"TerminalID"`
	RtnCode         string `json:"rtnCode"`
}

func (a *QueryResponse) GenerateCreditCheckParams() Request.Credit3dCheckParams {
	var data Request.Credit3dCheckParams
	data.TransDate = a.TransDate
	data.TransTime = a.TransTime
	data.ApproveCode = a.ApproveCode
	data.ResponseCode = a.ResponseCode
	data.ResponseMsg = a.ResponseMsg
	data.InstallType = a.InstallmentType
	data.FirstAmt = a.FirstAmt
	data.EachAmt = a.EachAmt
	data.Fee = a.Fee
	data.RedeemType = a.RedeemType
	data.RedeemUsed = a.RedeemUsed
	data.RedeemBalance = a.RedeemBalance
	data.CreditAmt = a.CreditAmt
	data.RiskMark = a.RiskMark
	data.FOREIGN = a.Foreign
	data.SECURE_STATUS = a.SecureStatus
	data.OrderID = a.OrderId
	data.TransAmt = a.TransAmt
	return data
}

func (a *GwCreditAuthData) GenerateGwCreditVoidData() GwCreditAuthData {
	var data GwCreditAuthData
	data.CardId = a.CardId
	data.MerchantId = a.MerchantId
	data.TerminalId = a.TerminalId
	data.OrderId = a.OrderId
	data.TramsAmount = a.TramsAmount
	data.TransTime = time.Now()
	data.TransType = Enum.CreditTransTypeVoid
	data.TransStatus = Enum.CreditTransStatusSuccess
	data.PayType = a.PayType
	data.CreditType = a.CreditType
	data.AuditStatus = Enum.CreditAuditRelease
	data.CaptureStatus = Enum.CreditCaptureInit
	return data
}

func (a *GwCreditAuthData) GenerateGwCreditRefundData() GwCreditAuthData {
	var data GwCreditAuthData
	data.CardId = a.CardId
	data.OrderId = a.OrderId
	data.MerchantId = a.MerchantId
	data.TerminalId = a.TerminalId
	data.TramsAmount = a.TramsAmount
	data.TransTime = a.TransTime
	data.TransType = Enum.CreditTransTypeRefund
	data.TransStatus = Enum.CreditTransStatusSuccess
	data.ApproveCode = a.ApproveCode
	data.PayType = a.PayType
	data.CreditType = a.CreditType
	data.AuditStatus = Enum.CreditAuditRelease
	data.CaptureStatus = Enum.CreditCaptureInit
	return data
}

type AuthResponse struct {
	Status string
	RtnURL string
}

type CancelRequest struct {
	OrderId    string
	MerchantId string
	TerminalId string
}

type QueryRequest struct {
	OrderId    string
	MerchantId string
	TerminalId string
}

type AuthResult struct {
	Installment     string `json:"Installment"`
	ResponseCode    string `json:"ResponseCode"`
	InstallmentType string `json:"InstallmentType"`
	RedeemType      string `json:"RedeemType"`
	Fee             string `json:"Fee"`
	ResponseMsg     string `json:"ResponseMsg"`
	CreditAmt       string `json:"CreditAmt"`
	MerchantID      string `json:"MerchantID"`
	FirstAmt        string `json:"FirstAmt"`
	RedeemUsed      string `json:"RedeemUsed"`
	PrivateData     string `json:"PrivateData"`
	OrderId         string `json:"OrderId"`
	TransMode       string `json:"TransMode"`
	EachAmt         string `json:"EachAmt"`
	TransAmt        string `json:"TransAmt"`
	ApproveCode     string `json:"ApproveCode"`
	TransDate       string `json:"TransDate"`
	TerminalID      string `json:"TerminalID"`
	TransTime       string `json:"TransTime"`
	RtnHtml         string `json:"rtnHtml"`
	RedeemBalance   string `json:"RedeemBalance"`
	RtnCode         string `json:"rtnCode"`
}

type AuthParams struct {
	MerchantID    string `form:"MerchantID"`    //特約商店代號
	TerminalID    string `form:"TerminalID"`    //端末機代號
	OrderID       string `form:"OrderID"`       //EC系統交易序號對應商店指派的「交易訂單編號」
	PAN           string `form:"PAN"`           //交易卡號(部份遮蓋)
	TransCode     string `form:"TransCode"`     //交易代碼(00 - 一般交易)
	TransDate     string `form:"TransDate"`     //交易日期(YYYYMMDD)
	TransTime     string `form:"TransTime"`     //交易時間(HHMMSS)
	TransAmt      string `form:"TransAmt"`      //交易金額
	ApproveCode   string `form:"ApproveCode"`   //授權碼
	ResponseCode  string `form:"ResponseCode"`  //回應碼
	ResponseMsg   string `form:"ResponseMsg"`   //回應訊息
	InstallType   string `form:"InstallType"`   //分期手續費計價方式(E – 外加，I – 內合)
	Install       int    `form:"Install"`       //分期期數
	FirstAmt      int    `form:"FirstAmt"`      //首期金額
	EachAmt       int    `form:"EachAmt"`       //每期金額
	Fee           int    `form:"Fee"`           //手續費
	RedeemType    string `form:"RedeemType"`    //紅利折抵方式(1 –全額，2 –部份)
	RedeemUsed    int    `form:"RedeemUsed"`    //紅利折抵點數
	RedeemBalance int    `form:"RedeemBalance"` //紅利餘額
	CreditAmt     int    `form:"CreditAmt"`     //持卡人自付額
	RiskMark      string `form:"RiskMark"`      //風險卡號
	FOREIGN       string `form:"FOREIGN"`       //國外卡
	SECURE_STATUS string `form:"SECURE_STATUS"` //3D 認證結果
}

//信用卡代碼表
type BankBinCode struct {
	Id         int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	BinNumber  string    `xorm:"varchar(10) notnull comment('銀行代碼')"`
	BankName   string    `xorm:"varchar(50) notnull comment('銀行名稱')"`
	CardType   string    `xorm:"varchar(10) notnull comment('卡別')"`
	IsDebit    int64     `xorm:"tinyint(1) default 0 comment('是否為DEBIT')"`
	Status     int64     `xorm:"tinyint(1) default 0 comment('狀態')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type MemberGwCreditAuth struct {
	GwCreditAuthData `xorm:"extends"`
	MemberCardData   `xorm:"extends"`
}

type OrderGwCreditAuth struct {
	Auth  GwCreditAuthData `xorm:"extends"`
	Order OrderData        `xorm:"extends"`
}
