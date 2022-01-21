package entity

import (
	"api/services/Enum"
	"api/services/VO/Response"
	"api/services/util/tools"
	"time"
)

type WithdrawData struct {
	WithdrawId     string    `xorm:"pk varchar(50) unique comment('提領單號')"`
	UserId         string    `xorm:"varchar(50) notnull comment('使用者')"`
	BankCode       string    `xorm:"varchar(10) notnull comment('銀行代碼')"`
	BankName       string    `xorm:"varchar(50) notnull comment('銀行名稱')"`
	BankAccount    string    `xorm:"varchar(30) notnull comment('銀行帳戶')"`
	WithdrawAmt    int64     `xorm:"int(10) notnull comment('提領金額')"`
	WithdrawFee    int64     `xorm:"int(10) notnull comment('提領手續費')"`
	TransId        string    `xorm:"varchar(8) notnull comment('交易代碼')"`
	BatchId        string    `xorm:"varchar(20) comment('批次代碼')"`
	ExportTime     time.Time `xorm:"datetime comment('轉出時間')"`
	WithdrawStatus string    `xorm:"varchar(20) notnull comment('提領狀態')"`
	WithdrawTime   time.Time `xorm:"datetime notnull comment('提領時間')"`
	ResponseStatus string    `xorm:"varchar(20) default 'INIT' comment('回覆狀態')"`
	ResponseCode   string    `xorm:"varchar(10) comment('回覆代碼')"`
	ResponseTime   time.Time `xorm:"datetime comment('回覆時間')"`
	WorkStaff      string    `xorm:"varchar(50) comment('操作人員')"`
	CreateTime     time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime     time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type BankCodeData struct {
	Id         int64  `xorm:"pk int(10) unique autoincr comment('序號')"`
	BankName   string `xorm:"varchar(50) notnull comment('銀行名稱')"`
	BankCode   string `xorm:"varchar(10) notnull comment('銀行代碼')"`
	BranchCode string `xorm:"varchar(10) notnull comment('銀行代碼包含分行碼')"`
	BankStatus string `xorm:"varchar(20) notnull comment('狀態')"`
}

type EmailVerifyData struct {
	Id           int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	UserId       string    `xorm:"varchar(50) notnull comment('使用者')"`
	StoreId      string    `xorm:"varchar(50) comment('收銀機ID')"`
	Email        string    `xorm:"varchar(50) notnull comment('Mail')"`
	VerifyCode   string    `xorm:"varchar(100) notnull comment('驗證碼')"`
	SendTime     time.Time `xorm:"datetime notnull comment('發送時間')"`
	ExpiredTime  time.Time `xorm:"datetime notnull comment('到期時間')"`
	VerifyStatus string    `xorm:"varchar(20) notnull comment('驗證狀態, WAIT, SUCCESS')"`
	VerifyType   string    `xorm:"varchar(10) notnull comment('驗證類型')"`
	CreateTime   time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime   time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type MemberWithdrawData struct {
	Id         int64     `xorm:"pk int(10) unique autoincr comment('序號')"`
	UserId     string    `xorm:"varchar(50) notnull comment('使用者')"`
	BankName   string    `xorm:"varchar(50) notnull comment('銀行名稱')"`
	BankCode   string    `xorm:"varchar(10) notnull comment('銀行代碼')"`
	Account    string    `xorm:"varchar(30) notnull comment('帳戶')"`
	ActDefault int64     `xorm:"tinyint(1) default 0 comment('預設帳戶')"`
	Status     string    `xorm:"varchar(20) notnull comment('狀態')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type DonateData struct {
	Id           int64  `xorm:"pk int(10) unique autoincr comment('序號')"`
	DonateName   string `xorm:"varchar(50) notnull comment('受捐贈機關或團體名稱')"`
	DonateCode   string `xorm:"varchar(10) notnull comment('捐贈碼')"`
	DonateShort  string `xorm:"varchar(20) comment('受捐贈機關或團體簡稱')"`
	DonateBan    string `xorm:"varchar(10) notnull comment('受捐贈機關或團體統編')"`
	DonateCity   string `xorm:"varchar(10) notnull comment('縣市')"`
	DonateStatus string `xorm:"varchar(10) notnull comment('狀態')"`
}

type QueryWithdraw struct {
	Withdraw WithdrawData `xorm:"extends"`
	Member   MemberData   `xorm:"extends"`
}

func (w *QueryWithdraw) GetWithdraw() Response.Withdraw {
	var data Response.Withdraw
	data.WithdrawDate = w.Withdraw.WithdrawTime.Format("2006/01/02 15:04")
	data.WithdrawId = w.Withdraw.WithdrawId
	data.WithdrawAmount = w.Withdraw.WithdrawAmt
	data.WithdrawFee = w.Withdraw.WithdrawFee
	data.BankAccount = tools.MaskerBankAccount(w.Withdraw.BankAccount)
	data.BankName = w.Withdraw.BankName
	data.Buyer = w.Member.Mphone //拿掉隱碼
	data.TerminalId = w.Member.TerminalId
	data.WithdrawType = tools.GetWithdrawType(w.Withdraw.BankCode)
	data.WithdrawStatus = w.Withdraw.ResponseStatus
	data.WithdrawStatusText = Enum.WithdrawStatus[w.Withdraw.ResponseStatus]
	return data
}

type IndustryData struct {
	IndustryId string    `xorm:"pk varchar(10) unique comment('行業代碼')"`
	Mcc        string    `xorm:"varchar(10) comment('MCC')"`
	Category   string    `xorm:"varchar(20) notnull comment('行業類別')"`
	Industry   string    `xorm:"varchar(50) notnull comment('行業')"`
	Sort       string    `xorm:"varchar(10) notnull comment('排序')"`
	CreateTime time.Time `xorm:"datetime notnull comment('建立時間')"`
}
