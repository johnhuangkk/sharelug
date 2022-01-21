package ExcelVo

type MemberReportVo struct {
	MemberId   string //會員帳號
	TerminalId string //會員代碼
	Balance    int64  //會員餘額
}

type OrderReportVo struct {
	Id          int64  //序號
	OrderTime   string //交易日期
	PayWay      string //付款方式
	PayWayTime  string //付款日期
	ShipType    string //出貨方式
	ShipTime    string //出貨日期
	OrderId     string //訂單編號
	SellerId    string //賣方會員代碼
	BuyerId     string //買方會員代碼
	Amount      int64  //交易金額
	PlatformFee int64  //平台費用
	IsFee       string //是否已付平台費用
}

type BankReportVo struct {
	Account      string //賣家帳號
	TerminalId   string //會員代碼
	Category     string //種類
	Name         string //姓名
	Identity     string //證號
	Head         string //負責人
	HeadIdentity string //負責人證號
	CompanyAddr  string //地址
	Store1       string //收銀機名稱1
	Store2       string //收銀機名稱2
	Store3       string //收銀機名稱3
	Store4       string //收銀機名稱4
	Store5       string //收銀機名稱5
}

type ShipReportVo struct {
	Id              int64
	OrderId         string
	ShipIdn         string
	ShipNumber      string
	ProductName     string
	Price           string
	Pieces          string
	ReceiverName    string
	ReceiverPhone   string
	ReceiverCode    string
	ReceiverAddress string
	OrderMemo       string
	BuyerNotes      string
}

type DayStatementVo struct {
	Id              int64  //序號
	TransactionDate string //交易日期
	TransactionType string //入金
	PaymentType     string //付款方式
	OrderId         string //訂單編號
	SellerId        string //賣方會員代碼
	BuyerId         string //買方會員代碼
	Amount          int64  //交易金額
	PlatformFee     int64  //平台費用
}

type InvoiceReportVo struct {
	InvoiceYm     string //發票年月
	InvoiceNumber string //發票號碼
	CreateTime    string //開立時間
	Products      string //購買內容
	Sales         int64  //應稅金額
	Tax           int64  //稅額
	Amount        int64  //發票金額
	Identifier    string //公司統編
	InvoiceStatus string //發票狀態
}

type UserInvoiceReportVo struct {
	OrderTime     string
	OrderId       string
	InvoiceNumber string
	CreateTime    string
	Amount        int64
}

type CouponUsedReportVo struct {
	Id             int64 //流水號
	Code           string
	Status         string
	OrderId        string
	TransTime      string
	BuyerPhone     string
	Amount         int64
	DiscountAmount int64
}

type SpecialStoreReportVo struct {
	MerchantId    string
	Terminal3DId  string
	TerminalN3DId string
	ChStoreName   string
	EnStoreName   string
}

type SpecialStoreRecordVo struct {
	Id             int64
	RepresentId    string
	RepresentFirst string
	RepresentLast  string
	Addr           string
	AddrEn         string
	CityCode       string
	CityName       string
	Represent      string
	MerchantId     string
	TerminalId     string
	Terminal3DId   string
	StoreName      string
	StoreNameEn    string
	JobCode        string
	MccCode        string
}
