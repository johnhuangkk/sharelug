package entity

import (
	"api/services/VO/Response"
	"time"
)

//資料結構
type MemberData struct {
	Uid              string    `xorm:"pk unique varchar(50) notnull comment('使用者ID')"`
	Mphone           string    `xorm:"varchar(10) notnull comment('手機號碼(帳號)')"`
	Email            string    `xorm:"varchar(50) comment('電子郵件')"`
	VerifyEmail      string    `xorm:"varchar(50) comment('驗證中電子郵件')"`
	VerifyIdentity   int64     `xorm:"tinyint(1) default 0 comment('是否身份認證')"`
	VerifyBusiness   int64     `xorm:"tinyint(1) default 0 comment('是否完成行業別')"`
	Identity         string    `xorm:"varchar(50) comment('身份證號')"`
	IdentityName     string    `xorm:"varchar(50) comment('身份認證使用姓名')"`
	Username         string    `xorm:"varchar(100) comment('使用者姓名')"`
	Picture          string    `xorm:"varchar(255) comment('大頭貼')"`
	PushToken        string    `xorm:"varchar(300) comment('push token')"`
	CertifiedStatus  string    `xorm:"varchar(10) comment('認證狀態')"`
	MemberStatus     string    `xorm:"varchar(10) comment('帳戶狀態')"`
	TerminalId       string    `xorm:"varchar(12) comment('會員代號')"`
	RealName         string    `xorm:"varchar(100) comment('訂購人姓名')"`
	SendName         string    `xorm:"varchar(20) comment('寄件人姓名')"`
	InvoiceCarrier   string    `xorm:"varchar(50) comment('發票載具號碼')"`
	Error            int64     `xorm:"int(2) notnull default 0 comment('登入錯誤次數')"`
	ErrorTime        time.Time `xorm:"datetime comment('登入錯誤時間')"`
	UpgradeType      string    `xorm:"varchar(10) notnull default 'RENEW' comment('加值狀態')"`
	UpgradeLevel     int64     `xorm:"int(10) notnull default 0 comment('加值等級')"`
	UpgradeExpire    time.Time `xorm:"datetime comment('加值到期時間')"`
	LastTime         time.Time `xorm:"datetime notnull comment('最後登入時間')"`
	RegisterTime     time.Time `xorm:"datetime notnull comment('註冊時間')"`
	Category         string    `xorm:"varchar(10) default 'MEMBER' comment('帳戶類別')"`
	CompanyName      string    `xorm:"varchar(50) comment('公司名稱')"`
	CompanyAddr      string    `xorm:"text comment('公司地址')"`
	CompanyAddressEn string    `xorm:"varchar(100) comment('公司地址英')"`
	Representative   string    `xorm:"varchar(50) comment('公司代理人/負責人')"`
	RepresentativeId string    `xorm:"varchar(50) comment('公司代理人/負責人身分證證號')"`
	RepresentLast    string    `xorm:"varchar(50) comment('公司代理人英性')"`
	RepresentFirst   string    `xorm:"varchar(50) comment('公司代理人英名')"`
	MccCode          string    `xorm:"varchar(10) comment('銀行行業代碼')"`
	CityCode         string    `xorm:"varchar(10) comment('區域代碼')"`
	CityNameEn       string    `xorm:"varchar(15) comment('區域英文名')"`
	JobCode          string    `xorm:"varchar(10) comment('行業代碼')"`
	Capital          float64   `xorm:"decimal(10,2) comment('資本額萬')"`
	ZipCode          string    `xorm:"varchar(5) comment('郵遞區號')"`
	Establish        string    `xorm:"varchar(10) comment('公司建立時間yyyymmdd')"`
	Contact          string    `xorm:"varchar(10) comment('聯絡人')"`
	ContactPhone     string    `xorm:"varchar(10) comment('聯絡人電話')"`
	ReportBank       bool      `xorm:"tinyint(1) default 0 comment('銀行申報賣家資料')"`
	Unsubscribe      bool      `xorm:"tinyint(1) default 0 comment('退訂電子報')"`
	CreateTime       time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime       time.Time `xorm:"datetime notnull comment('更新時間')"`
}

func (m *MemberData) GetMemberLoginInfo() Response.MemberInfo {
	var mi Response.MemberInfo
	mi.Uid = m.Uid
	mi.Email = m.Email
	mi.Mphone = m.Mphone
	mi.Username = m.Username
	mi.RealName = m.RealName
	mi.IdentityName = m.IdentityName
	mi.Picture = m.Picture
	mi.Category = m.Category
	mi.VerifyIdentity = m.VerifyIdentity
	mi.VerifyBusiness = m.VerifyBusiness
	mi.UpgradeLevel = m.UpgradeLevel
	return mi
}

type UserInfoResponse struct {
	MemberData MemberData    `xorm:"extends"`
	StoreData  StoreDataResp `xorm:"extends"`
}

type QueryUserStore struct {
	MemberData MemberData `xorm:"extends"`
	StoreData  StoreData  `xorm:"extends"`
}

type MemberCarrierData struct {
	MemberId    string    `xorm:"pk unique varchar(50) notnull comment('會員ID')"`
	InvoiceType string    `xorm:"varchar(10) notnull default 'PERSONAL' comment('發票設定')"`
	CompanyBan  string    `xorm:"varchar(10) comment('統一編號')"`
	CompanyName string    `xorm:"varchar(50) comment('公司名稱')"`
	DonateBan   string    `xorm:"varchar(10) comment('捐贈碼')"`
	CarrierType string    `xorm:"varchar(10) notnull default 'MEMBER' comment('發票載具')"`
	CarrierId   string    `xorm:"varchar(64) comment('載具ID')"`
	CreateTime  time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime  time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type MemberWithSpecialStore struct {
	MemberData      MemberData      `xorm:"extends"`
	KgiSpecialStore KgiSpecialStore `xorm:"extends"`
}

type MemberSendKgiBank struct {
	Id         int64     `xorm:"pk int(10) autoincr"`
	MerchantId string    `xorm:"varchar(30) comment('特店編號')"`
	Uid        string    `xorm:"varchar(50) notnull comment('使用者ID')"`
	IsSend     bool      `xorm:"tinyint(1) default 0 comment('銀行資料是否上傳')"`
	FileName   string    `xorm:"varchar(30) comment('上傳檔名')"`
	Created    time.Time `xorm:"timestamp created  comment('建立時間')"`
}

type MemberWithSendToSpecial struct {
	MemberSendKgiBank MemberSendKgiBank      `xorm:"extends"`
	Member            MemberWithSpecialStore `xorm:"extends"`
}
