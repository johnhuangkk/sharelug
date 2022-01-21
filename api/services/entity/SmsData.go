package entity

type SmsRequest struct {
	Username string `form:"username"` //使用者帳號 必要
	Password string `form:"password"` //使用者密碼 必要
	Dstaddr  string `form:"dstaddr"`  //收訊人之手機號碼，格式為:0912345678 必要
	Smbody   string `form:"smbody"`   //簡訊內容 此欄位進行URLEnocde 必要
	//destname string //收訊人名稱。
	//dlvtime  string //簡訊預約時間
	//vldtime  string //簡訊有效期限
	//response string //狀態主動回報網址
	//clientid string //客戶簡訊ID
	//objectID string //批次名稱
}

type SmsResponse struct {
	Clientid     string `form:"clientid"`
	Msgid        string `form:"msgid"`
	Statuscode   string `form:"statuscode"`
	AccountPoint string `form:"AccountPoint"`
	Duplicate    string `form:"Duplicate"`
}

type SmsSubmitReq struct {
	SysId         string `xml:"SysId" `
	SrcAddress    string `xml:"SrcAddress"`
	DestAddress   string `xml:"DestAddress"`
	SmsBody       string `xml:"SmsBody"`
	DrFlag        bool   `xml:"DrFlag"`
	FirstFailFlag bool   `xml:"FirstFailFlag"`
}

type SmsMultiSubmitReq struct {
	SysId         string   `xml:"SysId" `
	SrcAddress    string   `xml:"SrcAddress"`
	DestAddress   []string `xml:"DestAddress"`
	SmsBody       string   `xml:"SmsBody"`
	DrFlag        bool     `xml:"DrFlag"`
	FirstFailFlag bool     `xml:"FirstFailFlag"`
}

type SmsResult struct {
	ResultCode string `xml:"ResultCode"`
	ResultText string `xml:"ResultText"`
	MessageId  string `xml:"MessageId"`
}
