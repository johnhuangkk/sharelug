package AtmResponse

type SMX  struct {
	Header Header `xml:"Header"`
	SvcRq  SvcRq  `xml:"SvcRq"`
}

type Header struct {
	TxnId	string	`xml:"TxnId"`			//交易代碼
	Status	Status	`xml:"Status"`			//交易結果
}

type Status struct {
	StatusCode	string 	`xml:"StatusCode"`	//交易結果代碼(0:成功 其它:錯誤)
	StatusDesc	string	`xml:"StatusDesc"`	//交易結果錯誤代碼說明
}

type SvcRq struct {
	ACCTNO		string 	`xml:"ACCTNO"` 		//實際入帳帳號(固定 16 碼不足右靠左補 0)
	INQOPTNO	string	`xml:"INQOPTNO"`	//查詢委託單位代碼(固定 5 碼不足右靠左補 0，00000 表示全部)
	BDATE		string	`xml:"BDATE"`		//查詢帳務日起日 YYYYMMDD
	EDATE		string	`xml:"EDATE"`		//查詢帳務日迄日 YYYYMMDD
	SBDATE		string	`xml:"SBDATE"`		//查詢交易日起日 YYYYMMDD
	SEDATE		string	`xml:"SEDATE"`		//查詢交易日迄日 YYYYMMDD
	SBTIME		string	`xml:"SBTIME"`		//查詢交易時間起日 HHMM
	SETIME		string	`xml:"SETIME"`		//查詢交易時間迄日 HHMM
	VACNO		string	`xml:"VACNO"`		//查詢特定虛擬帳號
	TEMPDATA	string	`xml:"TEMPDATA"`	//中間鍵值資料。(固定 5 碼) 作為下一筆資料讀取用。如為空白為最後一筆
	DETAIL		DETAIL	`xml:"DETAIL"`
}

type DETAIL struct {
	ITEMNO		string	`xml:"ITEMNO"`		//序號
	TXNDATE		string	`xml:"TXNDATE"`		//交易日 YYYYMMDD
	TXNTIME		string	`xml:"TXNTIME"`		//交易時間 HHmmSS
	REODATE		string	`xml:"REODATE"`		//帳務日 YYYYMMDD
	REOTYPE		string	`xml:"REOTYPE"`		//代收通路中文說明
	SOURCE		string	`xml:"SOURCE"`		//代收通路代碼 (A:臨櫃 B/P:語音 C:網銀 D:行動銀 E/R:匯款 F:FXML G:全繳 J:ADM M:MOD T:ATM X:eATM)
	REOMAIL		string	`xml:"REOMAIL"`		//虛擬帳號(固定 16 碼不足右靠左補 0)
	CURRENCY	string	`xml:"CURRENCY"`	//幣別
	REODCCODE	string	`xml:"REODCCODE"`	//金額正負號(+/-;負號為沖正交易)
	AMT			string	`xml:"AMT"`			//金額 (交易幣別為外幣時含小數點及小數位數 2 位)
	REOSND		string	`xml:"REOSND"`		//代收行/轉出行/匯款行 外幣交易時為 SWIFT CODE(自行轉帳交易固定為 CDIBTWTPXXX)
	REONAME		string	`xml:"REONAME"`		//存繳人/匯款人/附言 轉出帳號(固定 16 碼不足右靠左補 0)
	REOIDNO		string	`xml:"REOIDNO"`		//外幣轉帳交易之轉出帳號對映統編/身份證字號 台幣交易及外幣匯入匯款時無此資訊
	NUM			string	`xml:"NUM"`			//交易序號(為每一帳務日之 KEY 值;若需判斷入金交易是否重覆， 亦請加上帳務日一起判斷。)
}
