package entity

import "time"

type TransferData struct {
	Id              int       `xorm:"pk int(10) unique autoincr comment('序號')"`
	OrderId         string    `xorm:"varchar(50) notnull comment('訂單編號')"`
	BankAccount     string    `xorm:"varchar(50) notnull comment('轉帳帳號')"`
	BankCode        string    `xorm:"varchar(50) notnull comment('銀行代碼')"`
	BankName        string    `xorm:"varchar(20) notnull comment('銀行名稱')"`
	Amount          int64     `xorm:"varchar(50) notnull comment('轉帳銀額')"`
	Currency        string    `xorm:"varchar(50) notnull comment('幣別')"`
	TransType       string    `xorm:"varchar(20) notnull comment('交易類別')"`
	ExpireDate      time.Time `xorm:"datetime comment('到期日')"`
	RecdBankAccount string    `xorm:"varchar(50) comment('轉入帳號')"`
	RecdAmount      string    `xorm:"varchar(50) comment('轉入金額')"`
	RecdDate        time.Time `xorm:"datetime comment('轉入時間')"`
	Seqno           string    `xorm:"varchar(50) comment('交易序號')"`
	TransferStatus  string    `xorm:"varchar(50) notnull comment('狀態')"`
	CreateTime      time.Time `xorm:"datetime notnull comment('建立時間')"`
	UpdateTime      time.Time `xorm:"datetime notnull comment('更新時間')"`
}

type TransferParams struct {
	Resend    string `json:"RESEND"`     //訊息種類(I=首次 R=重送;本服務無重送功能)
	Head      string `json:"HEAD"`       //訊息代碼(頭尾 tag，無須理會)
	Tdate     string `json:"TDATE"`      //交易日期 YYYYMMDD
	Ttime     string `json:"TTIME"`      //交易時間 HHMMSS
	Date      string `json:"DATE"`       //帳務日期 YYYYMMDD
	Accno     string `json:"ACCNO"`      //虛擬帳號(固定 16 碼不足右靠左補 0)
	DepType   string `json:"DepType"`    //借貸別皆為貸(1=借,2=貸)
	Currency  string `json:"CURRENCY"`   //交易幣別
	Sign      string `json:"SIGN"`       //金額正負號(+/-;負號為沖正交易)
	Amt       string `json:"AMT"`        //交易金額 (交易幣別為外幣時含小數點及小數位數 2 位)
	Type      string `json:"TYPE"`       //交易種類 (A:臨櫃 B/P:語音 C:網銀 D:行動銀 E/R:匯款 F:FXML G:eBill J:ADM M:MOD T:ATM X:eATM 0:其它)
	Raccno    string `json:"RACCNO"`     //轉出帳號(銀行代號 3 位+帳號 16 位)
	SwiftCode string `json:"SWIFT CODE"` //外幣轉出行/匯款行 SWIFT CODE(自行轉帳交易固定為 CDIBTWTPXXX)
	Note      string `json:"NOTE"`       //備註(左靠右補空白)
	Atype     string `json:"ATYPE"`      //帳戶類別(預留帶空白)
	Idno      string `json:"IDNO"`       //外幣轉帳交易之轉出帳號對映統編/身份證字號 台幣交易及外幣匯入匯款時無此資訊
	Baccno    string `json:"BACCNO"`     //實際入款帳號(固定 16 碼不足右靠左補 0)
	Seqno     string `json:"SEQNO"`      //交易序號
	Eor       string `json:"EOR"`        //訊息代碼(頭尾 tag，無須理會)
}

type SMX struct {
	Header Header `xml:"Header"`
	SvcRq  SvcRq  `xml:"SvcRq"`
}

type Header struct {
	Date     string `xml:"Date"`     //傳輸日期 YYYYMMDD
	Time     string `xml:"Time"`     //傳輸時間 HHmmSS
	SenderID string `xml:"SenderID"` //傳送端公司統編。必要欄位
	ReceID   string `xml:"ReceID"`   //接收端代號(固定 809)
	Password string `xml:"Password"` //密碼(銀行端提供)。必要欄位
	TxnId    string `xml:"TxnId"`    //交易代碼(固定 V522-可查詢台/外幣交易)。必要欄位
}

type SvcRq struct {
	IDNO     string `xml:"IDNO"`     //客戶統編。必要欄位
	ACNO     string `xml:"ACNO"`     //實際入帳帳號。必要欄位
	INQOPTNO string `xml:"INQOPTNO"` //特定委託單位代碼：固定 5 碼不足右靠左補 0
	BDATE    string `xml:"BDATE"`    //查詢帳務日起日 YYYYMMDD。與交易日二擇一
	EDATE    string `xml:"EDATE"`    //查詢帳務日迄日 YYYYMMDD。與交易日二擇一
	SBDATE   string `xml:"SBDATE"`   //查詢交易日起日 YYYYMMDD。與帳務日二擇一
	SEDATE   string `xml:"SEDATE"`   //查詢交易日迄日 YYYYMMDD。與帳務日二擇一
	SBTIME   string `xml:"SBTIME"`   //查詢交易時間起日 HHMM。可空白，若有值限以交易日查詢
	SETIME   string `xml:"SETIME"`   //查詢交易時間迄日 HHMM。可空白，若有值限以交易日查詢
	VACNO    string `xml:"VACNO"`    //查詢特定虛擬帳號。 限以帳務日查詢且查詢區間限 7 日內
	TEMPD    string `xml:"TEMPD"`    //中間鍵值資料。必要欄位 第1次為空白。
}

type OrderAndTransfer struct {
	Transfer TransferData `xorm:"extends"`
	Order    OrderData    `xorm:"extends"`
}
