package KgiBank

type Header struct {
	Flag       string //檔頭旗標 固定為 "H"
	MerchantId string //特店代號
	SendDate   string //送檔日期
	Seq        string //序號
	Count      string //總筆數
	Symbol     string //金額正負號
	ToTal      string //總金額
	Filler     string //FILLER填空白
}

type Body struct {
	MerchantId   string //特店代號
	TerminalId   string //端末機代號
	OrderId      string //訂單編號
	Space        string //空白
	TranAmount   string //交易金額
	AuthCode     string //授權碼
	TranType     string //交易碼
	TranDate     string //交易日期
	Custom       string //使用者自訂欄位
	CardInfo     string //持卡人資訊
	ProcessDate  string //帳單處理日期
	ResponseCode string //回應碼
	ResponseMsg  string //回應訊息
	BatchSeq     string //Batch and seq. No.
	Mark         string //分期付款或紅利積點註記
	NumberOfPay  string //分期數
	FirstPayment string //首期金額
	EachPayment  string //每期金額
	Fees         string //手續費
	DeductPoint  string //本次扣底點數
	Symbol       string //金額正負號
	PointBalance string //持卡人點數餘額
	Deductible   string //持卡人自付金額
	PaymentDate  string //付款日
	Verify       string //認證結果
	Foreign      string //國外卡
	Reserved     string //預留
}