package AtmRequest

type SMX  struct {
	Header Header `xml:"Header"`
	SvcRq  SvcRq  `xml:"SvcRq"`
}

type Header struct {
	Date		string	`xml:"Date"`		//傳輸日期 YYYYMMDD
	Time		string	`xml:"Time"`		//傳輸時間 HHmmSS
	SenderID	string	`xml:"SenderId"`	//傳送端公司統編。必要欄位
	ReceID		string	`xml:"ReceId"`		//接收端代號(固定 809)
	Password	string	`xml:"Password"`	//密碼(銀行端提供)。必要欄位
	TxnId		string	`xml:"TxnId"`		//交易代碼(固定 V522-可查詢台/外幣交易)。必要欄位
}

type SvcRq struct {
	IDNO		string 	`xml:"IDNO"`		//客戶統編。必要欄位
	ACNO		string	`xml:"ACNO"`		//實際入帳帳號。必要欄位
	INQOPTNO	string	`xml:"INQOPTNO"`	//查詢委託單位代碼。必要欄位 0:全部
	BDATE		string	`xml:"BDATE"`		//查詢帳務日起日 YYYYMMDD。與交易日二擇一
	EDATE		string	`xml:"EDATE"`		//查詢帳務日迄日 YYYYMMDD。與交易日二擇一
	SBDATE		string	`xml:"SBDATE"`		//查詢交易日起日 YYYYMMDD。與帳務日二擇一
	SEDATE		string	`xml:"SEDATE"`		//查詢交易日迄日 YYYYMMDD。與帳務日二擇一
	SBTIME		string	`xml:"SBTIME"`		//查詢交易時間起日 HHMM。可空白，若有值限以交易日查詢
	SETIME		string	`xml:"SETIME"`		//查詢交易時間迄日 HHMM。可空白，若有值限以交易日查詢
	VACNO		string	`xml:"VACNO"`		//查詢特定虛擬帳號。限以帳務日查詢且查詢區間限 7 日內 委託單位代號不得空白或為 0
	TEMPD		string	`xml:"TEMPD"`		//中間鍵值資料。必要欄位 第 1 次為空白。 第 2 次開始直接由前一筆下行資料中遞送(需將此欄位填入上次 下行電文的 <TEMPDATA>;直至回傳<TEMPDATA>值為空白)
}


func Request() {

}
