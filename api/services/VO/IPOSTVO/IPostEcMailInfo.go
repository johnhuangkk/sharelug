package IPOSTVO

/**
電商平台上傳i郵箱郵件交寄資料處理訊息之交換方式及規格
*/
type UP_ECMAILINFO struct {
	TimeStamp string
	Token     string
	PayLoad   PayLoad
}

type PayLoad struct {
	APLBR_VipNO   string
	APLBR_VipName string
	MailNo        string
	// v 可用電話後四碼
	PrnAuthCode string
	EcOrderNo   string
	EcOrderDate string
	// v 選填
	PaymentAct string
	// v 選填
	PaymentAmt  int
	PrdName     string
	SenderName  string
	SenderPhone string
	// v 001:一般地址 002:IBOX 地址 003:郵政信箱 004:特種郵政信箱
	SenderAddrType string
	SenderZipCode  string
	SenderAddr     string
	// v 選填
	SenderAddrIBoxId string
	ReceiverName     string
	ReceiverPhone    string
	ReceiverZipCode  string
	// v 001:一般地址 002:IBOX 地址 003:郵政信箱 004:特種郵政信箱
	ReceiverAddrType string
	ReceiverAddr     string
	// v 選填
	ReceiverAddrIBoxID string
	// v 1 : 退回 2 : 拋棄
	ReturnType int
	// v 收件有效截止日期 (含當日) yyyyMMdd
	ValidDate string
	// v 列印及前台寄件有 效截止日期(含當日) yyyyMMdd
	ValidDate_PrintAndSend string
	Remark                 string
}

// 郵箱郵遞區號和完整地址
type IPostZipAddress struct {
	Id      string
	Alias   string
	Zip     string
	Address string
	Status  string
}

// 郵箱取號設定賣家地址
type SellerAddress struct {
	Id string
	Zip string
	Address string
}

/**
電商平台上傳i郵箱郵件交寄資料處理訊息 回應
*/
type RspEcMailInfo struct {
	MailNo, Failures string
}
