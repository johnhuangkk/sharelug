package Request

type SearchOrderRequest struct {
	Search OrderRequest `json:"Search"`
	//Limit  int          `json:"Limit"`
	//Start  int          `json:"Start"`
}

type OrderRequest struct {
	OrderDate     string `json:"OrderDate"`     //訂購區間
	PaymentDate   string `json:"PaymentDate"`   //付款區間
	ShipDate      string `json:"ShipDate"`      //出貨區間
	Seller        string `json:"Seller"`        //會員帳號
	Buyer         string `json:"Buyer"`         //訂購帳號
	Receiver      string `json:"Receiver"`      //收件人
	PaymentType   string `json:"PaymentType"`   //付款方式
	ShipType      string `json:"ShipType"`      //出貨方式
	ReturnId      string `json:"ReturnId"`      //退貨編號
	SellerId      string `json:"SellerId"`      //會員代碼
	BuyerName     string `json:"BuyerName"`     //訂購人
	ReceiverPhone string `json:"ReceiverPhone"` //收件電話
	OrderAmount   int64  `json:"OrderAmount"`   //訂單金額
	ShipNumber    string `json:"ShipNumber"`    //出貨單號
	ShipStatus    string `json:"ShipStatus"`    //退貨狀態
	StoreName     string `json:"StoreName"`     //賣場名稱
	OrderIp       string `json:"OrderIp"`       //訂購IP
	ReceiverAddr1 string `json:"ReceiverAddr1"` //收貨地址
	ProductAmount int64  `json:"ProductAmount"` //商品金額
	ShipMode      string `json:"ShipMode"`      //物流業者
	RefundId      string `json:"RefundId"`      //退款編號
	OrderId       string `json:"OrderId"`       //訂單編號
	OrderStatus   string `json:"OrderStatus"`   //訂單狀態
	ProductId     string `json:"ProductId"`     //商品編號
	ShipFee       int64  `json:"ShipFee"`       //運費金額
	ReceiverAddr2 string `json:"ReceiverAddr2"` //收件地址
	RefundStatus  string `json:"RefundStatus"`  //退款狀態
}

type SearchMemberRequest struct {
	Search MemberRequest `json:"Search"`
}

type MemberRequest struct {
	Account    string `json:"Account"`    //會員帳號
	Nickname   string `json:"Nickname"`   //會員暱稱
	StoreName  string `json:"StoreName"`  //賣場名稱
	TerminalId string `json:"TerminalId"` //會員代碼
}

type ErpSearchOrderRequest struct {
	OrderId string `json:"OrderId"`
}

type ErpRequest struct {
	AuditType string   `json:"AuditType"`
	OrderId   []string `json:"OrderId"`
}

type ErpAuditMemoRequest struct {
	OrderId string `json:"OrderId"`
	Memo    string `json:"Memo"`
}

type ErpAuditListRequest struct {
	Tabs string `json:"Tabs"`
}

type ErpDemoteRequest struct {
	UserId string `json:"UserId"`
	Level  int64  `json:"Level"`
}

type SearchWithdrawRequest struct {
	Tabs   string         `json:"Tabs"`
	Search SearchWithdraw `json:"Search"`
}

type ChangeWithdrawRequest struct {
	Status     string   `json:"Status"`
	WithdrawId []string `json:"WithdrawId"`
}

type SearchWithdraw struct {
	WithdrawDate   string `json:"WithdrawDate"`   //提領區間
	Buyer          string `json:"Buyer"`          //會員帳號
	WithdrawStatus string `json:"WithdrawStatus"` //提領狀態
	WithdrawId     string `json:"WithdrawId"`     //提領編號
	BuyerEmail     string `json:"BuyerEmail"`     //會員Email
}

type CvsSendCheckedRequest struct {
	Duration string `json:"Duration"`
	Checked  string `json:"Checked,omitempty"`
	Type     string `json:"Type"`
}

type ErpSearchProductRequest struct {
	Tab           string `json:"Tab"`
	UserAccount   string `json:"Account"`       //會員帳號
	StoreName     string `json:"StoreName"`     //收銀機名稱
	UserId        string `json:"UserId"`        //會員代碼
	ProductStatus string `json:"ProductStatus"` //上架狀態
	CreateDate    string `json:"CreateDate"`    //上架區間 開立區間
	UpdateDate    string `json:"UpdateDate"`    //修改區間          －－－
	Amount        string `json:"Amount"`        //單價區間          金額區間
	ProductName   string `json:"ProductName"`   //商品名稱
	ProductId     string `json:"ProductId"`     //商品編號          帳單編號
	ShipMode      string `json:"ShipMode"`      //出貨方式
	PaymentMode   string `json:"PaymentMode"`   //付款方式
}

type ErpCustomerContact struct {
	OrderId string `json:"orderID"`
}

type CustomerReplyRequest struct {
	QuestionId string `json:"QuestionId"`
	Content    string `json:"Content"`
}

type PlatformMessageRequest struct {
	Message string `json:"Message"`
	Type    string `json:"Type"`
}

type PlatformSendEdmMessageRequest struct {
	Message string `json:"Message"`
	Rule    string `json:"Rule"`
	Test    bool   `json:"Test"`
}

type GenerateShortRequest struct {
	Uri string `json:"Uri"`
}

type InvoiceTrackRequest struct {
	InvoiceTrack  string `json:"InvoiceTrack"`
	InvoiceBegin  string `json:"InvoiceBegin"`
	InvoiceEnd    string `json:"InvoiceEnd"`
	InvoicePeriod string `json:"InvoicePeriod"`
}

type SearchBalanceRequest struct {
	Account   string `json:"Account"`
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
}
