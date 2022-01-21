package Response

type AuditListResponse struct {
	AuditTabs     AuditTabs       `json:"AuditTabs"`
	AuditListData []AuditListData `json:"AuditListData"`
}

type AuditTabs struct {
	Wait    int64 `json:"Wait"`
	Note    int64 `json:"Note"`
	Pending int64 `json:"Pending"`
	Capture int64 `json:"Capture"`
	Void    int64 `json:"Void"`
	Refund  int64 `json:"Refund"`
}

type AuditListData struct {
	PaymentTime  string `json:"PaymentTime"`  //付款日期
	OrderId      string `json:"OrderId"`      //訂單編號
	OrderStatus  string `json:"OrderStatus"`  //訂單狀態
	BuyerAccount string `json:"BuyerAccount"` //訂購人帳號
	BuyerName    string `json:"BuyerName"`    //訂購人姓名
	CardBank     string `json:"CardBank"`     //發卡銀行
	CardAccount  string `json:"CardAccount"`  //卡號
	OrderAmount  int64  `json:"OrderAmount"`  //交易金額
	CardVerify   string `json:"CardVerify"`   //3D認證
}

type GetCreditAuthResponse struct {
	BankName      string `json:"BankName"`      //收單銀行
	CardType      string `json:"CardType"`      //卡別
	ApproveCode   string `json:"ApproveCode"`   //授權碼
	ResponseCode  string `json:"ResponseCode"`  //授權回應碼
	ResponseMsg   string `json:"ResponseMsg"`   //授權回應訊息
	FirstTrans    string `json:"FirstTrans"`    //首次交易
	Foreign       string `json:"Foreign"`       //國外卡
	RiskNote      string `json:"RiskNote"`      //風險卡號註記
	RiskList      string `json:"RiskList"`      //風險名單
	IP            string `json:"Ip"`            //IP
	AuditStatus   string `json:"AuditStatus"`   //審單狀態
	PendingDate   string `json:"PendingDate"`   //待決日期
	NoteDate      string `json:"NoteDate"`      //照會日期
	RefusedDate   string `json:"RefusedDate"`   //拒絕日期
	ReleaseDate   string `json:"ReleaseDate"`   //放行日期
	CaptureStatus string `json:"CaptureStatus"` //請款狀態
	BatchDate     string `json:"CaptureDate"`   //請款日期
	VoidDate      string `json:"VoidDate"`      //取消授權日期
	RetreatDate   string `json:"RetreatDate"`   //退刷日期
	CaptureDate   string `json:"ResponseDate"`  //下載結果檔日期
}

type ErpOrderResponse struct {
	SellerName      string           `json:"SellerName"`      //賣家
	SellerAcct      string           `json:"SellerAcct"`      //賣家帳號
	StoreName       string           `json:"StoreName"`       //賣場名稱
	OrderDate       string           `json:"OrderDate"`       //訂購日期
	OrderId         string           `json:"OrderId"`         //訂單編號
	ProductTotal    int64            `json:"ProductTotal"`    //商品總金額
	OrderTotal      int64            `json:"OrderTotal"`      //訂單總金額
	ShippingFee     int64            `json:"ShippingFee"`     //運費
	PlatformFee     int64            `json:"PlatformFee"`     //平台服務費
	PlatformPayFee  int64            `json:"PlatformPayFee"`  //金流處理費
	PlatformShipFee int64            `json:"PlatformShipFee"` //物流處理費
	PlatformInfoFee int64            `json:"PlatformInfoFee"` //資訊處理費
	OrderDetail     []ErpOrderDetail `json:"OrderDetail"`     //訂單內容
}

type ErpOrderDetail struct {
	ShippingFee int64  `json:"ShippingFee"` //運費
	ProductName string `json:"ProductName"` //訂購商品名稱
	ProductSpec string `json:"ProductSpec"` //規格
	ProductAmt  int64  `json:"ProductAmt"`  //單價
	ProductQty  int64  `json:"ProductQty"`  //數量
	ShipMerge   bool   `json:"ShipMerge"`   //是否合併運費
}

type ErpOrderShipInfoResponse struct {
	ShipMethod     string `json:"ShipMethod"`     //出貨方式
	RecipientName  string `json:"Recipient"`      //收件人
	RecipientPhone string `json:"RecipientPhone"` //收件人電話
	ShipTrader     string `json:"ShipTrader"`     //物流業者
	ReceiptAddress string `json:"ReceiptAddress"` //收貨地點(地址/門市/據點)
	RefundAddress  string `json:"RefundAddress"`  //退回地址(未取、無法投遞)
}

type ErpOrderCustomerResponse struct {
	Contents    string `json:"Contents"`    //客服往來
	Remark      string `json:"Remark"`      //客服備註
	Messages    string `json:"Messages"`    //買賣方留言
	AuditRemark string `json:"AuditRemark"` //審單備註
}

type ErpOrderRefundResponse struct {
	ApplyDate     string `json:"ApplyDate"`     //申請退款日期
	RefundId      string `json:"RefundId"`      //退款編號
	Reason        string `json:"Reason"`        //退款原因
	SellerAccount string `json:"SellerAct"`     //賣家帳號
	SellerBalance int64  `json:"SellerBalance"` //賣家可提領餘額
	OrderAmount   int64  `json:"OrderAmount"`   //訂單金額
	RefundAmount  int64  `json:"RefundAmount"`  //退款金額
	BuyerAccount  string `json:"BuyerAccount"`  //買家帳號
	RefundStatus  string `json:"RefundStatus"`  //撥付狀態
}

type ErpOrderRefundListResponse struct {
	RefundTime   string `json:"RefundTime"`   //退款日期
	RefundId     string `json:"RefundId"`     //退款編號
	RefundReason string `json:"RefundReason"` //退款原因
	RefundAmount int64  `json:"RefundAmount"` //退款金額
	RefundStatus string `json:"RefundStatus"` //退款狀態
}

type ErpOrderReturnListResponse struct {
	ReturnTime    string `json:"ReturnTime"`    //退貨日期
	ReturnId      string `json:"ReturnId"`      //退貨編號
	ReturnStatus  string `json:"ReturnStatus"`  //退貨狀態
	ReturnProduct string `json:"ReturnProduct"` //退貨商品
	ReturnSpec    string `json:"ReturnSpec"`    //規格
	ReturnPrice   int64  `json:"ReturnPrice"`   //單價
	ReturnQty     int64  `json:"ReturnQty"`     //退貨數量
}

type SearchOrderResponse struct {
	Count  int64            `json:"Count"`
	Orders []OrdersResponse `json:"Orders"`
}

type OrdersResponse struct {
	OrderDate         string `json:"OrderDate"`         //訂購日期
	Seller            string `json:"Seller"`            //賣家帳號
	SellerId          string `json:"SellerId"`          //賣家會員代碼
	StoreName         string `json:"StoreName"`         //賣場名稱
	OrderId           string `json:"OrderId"`           //訂單編號
	OrderStatus       string `json:"OrderStatus"`       //訂單狀態
	OrderStatusText   string `json:"OrderStatusText"`   //
	OrderAmount       int64  `json:"OrderAmount"`       //訂單金額
	BuyerId           string `json:"BuyerId"`           //買家會員代碼
	PaymentType       string `json:"PaymentType"`       //付款方式
	PaymentTypeText   string `json:"PaymentTypeText"`   //
	PaymentTime       string `json:"PaymentTime"`       //付款時間
	ShipType          string `json:"ShipType"`          //出貨方式
	ShipTypeText      string `json:"ShipTypeText"`      //
	CaptureStatus     string `json:"CaptureStatus"`     //撥付狀態
	CaptureStatusText string `json:"CaptureStatusText"` //
	ProductAmount     int64  `json:"ProductAmount"`     //商品金額
	ShipFee           int64  `json:"ShipFee"`           //運費
	PlatformFee       int64  `json:"PlatformFee"`       //平台服務費
	Coupon            string `json:"Coupon"`
	CouponAmount      int64  `json:"CouponAmount"`
}

type SearchOrderDetailResponse struct {
	SellerName      string              `json:"SellerName"`  //賣家
	Seller          string              `json:"Seller"`      //賣家帳號
	StoreName       string              `json:"StoreName"`   //賣場
	OrderDate       string              `json:"OrderDate"`   //訂購日期
	OrderId         string              `json:"OrderId"`     //訂單編號
	OrderStatus     string              `json:"OrderStatus"` //訂單狀態
	OrderStatusText string              `json:"OrderStatusText"`
	OrderAmount     int64               `json:"OrderAmount"`     //訂單總金額
	ProductAmount   int64               `json:"ProductAmount"`   //商品總金額
	ShipFee         int64               `json:"ShipFee"`         //運費
	PlatformFee     int64               `json:"PlatformFee"`     //平台服務費
	PlatformPayFee  int64               `json:"PlatformPayFee"`  //金流處理費
	PlatformShipFee int64               `json:"PlatformShipFee"` //物流處理費
	PlatformInfoFee int64               `json:"PlatformInfoFee"` //資訊處理費
	Coupon          string              `json:"Coupon"`          //折扣碼
	CouponAmount    int64               `json:"CouponAmount"`    //折扣金額
	BuyerNotes      string              `json:"BuyerNotes"`      //買家留言
	SellerNotes     string              `json:"SellerNotes"`     //賣家留言
	Detail          []SearchOrderDetail `json:"Detail"`
}

type SearchOrderDetail struct {
	ShipFee       int64  `json:"ShipFee"`       //運費
	ProductName   string `json:"ProductName"`   //訂購商品
	ProductSpec   string `json:"ProductSpec"`   //規格
	Amount        int64  `json:"Amount"`        //單價
	Quantity      int64  `json:"Quantity"`      //數量
	PendingReturn int64  `json:"PendingReturn"` //待退貨數量
	Returned      int64  `json:"Returned"`      //已退貨數量
	IsMerge       int64  `json:"IsMerge"`       //是否合拼
}

type SearchOrderRefundResponse struct {
	Refund []RefundResponse `json:"Refund"`
	Return []ReturnResponse `json:"Return"`
}

type RefundResponse struct {
	RefundDate       string `json:"RefundDate"`       //退款日期
	RefundId         string `json:"RefundId"`         //退款編號
	Reason           string `json:"Reason"`           //退款原因
	Amount           int64  `json:"Amount"`           //退款金額
	RefundStatus     string `json:"RefundStatus"`     //退款狀態
	RefundStatusText string `json:"RefundStatusText"` //退款狀態
}

type ReturnResponse struct {
	ReturnDate       string `json:"ReturnDate"`       //退貨日期
	ReturnId         string `json:"ReturnId"`         //退貨編號
	ReturnStatus     string `json:"ReturnStatus"`     //退貨狀態
	ReturnStatusText string `json:"ReturnStatusText"` //退貨狀態
	ProductName      string `json:"ProductName"`      //退貨商品
	ProductSpec      string `json:"ProductSpec"`      //規格
	Amount           int64  `json:"Amount"`           //單價
	Quantity         int64  `json:"Quantity"`         //退貨數量
}

type ShippingResponse struct {
	ShipMode      string `json:"ShipMode"`      //出貨方式
	ShipNumber    string `json:"ShipNumber"`    //出貨單號
	Receiver      string `json:"Receiver"`      //收件人
	ReceiverPhone string `json:"ReceiverPhone"` //收件人電話
	ShipTime      string `json:"ShipTime"`      //出貨日期
	ShipTrader    string `json:"Operator"`      //物流業者
	ReceiverAddr  string `json:"ReceiverAddr"`  //收貨地點(地址/門市/據點)
	SwitchPieces  string `json:"SwitchPieces"`  //閉轉-取件店/點
	SellerAddr    string `json:"SellerAddr"`    //退回地址(未取、無法投遞)
	SwitchStore   string `json:"SwitchStore"`   //閉轉-退回店/點
}

type PaymentInfoResponse struct {
	BuyerName     string `json:"BuyerName"`     //訂購人
	BuyerPhone    string `json:"BuyerPhone"`    //訂購人電話
	BuyerId       string `json:"BuyerId"`       //會員代碼
	Receiver      string `json:"Receiver"`      //收件人
	ReceiverPhone string `json:"ReceiverPhone"` //收件人電話
	PaymentMode   string `json:"PaymentMode"`   //付款方式
	PaymentDate   string `json:"PaymentDate"`   //付款日期
	VoidDate      string `json:"VoidDate"`      //取消授權日期
	RefundDate    string `json:"RefundDate"`    //退刷日期
	BankName      string `json:"BankName"`      //發卡/匯出 - 銀行
	LastFour      string `json:"LastFour"`      //卡號/帳號 - 末四碼
	Amount        int64  `json:"Amount"`        //交易金額
	AcquirerBank  string `json:"AcquirerBank"`  //收單/收款銀行
}

type SearchWithdrawResponse struct {
	WithdrawTabs WithdrawTabs `json:"WithdrawTabs"`
	WithdrawList []Withdraw   `json:"WithdrawList"`
}

type WithdrawTabs struct {
	Wait    int64 `json:"Wait"`
	Pending int64 `json:"Pending"`
}

type Withdraw struct {
	WithdrawDate       string `json:"WithdrawDate"`       //申請提領日期
	WithdrawId         string `json:"WithdrawId"`         //提領編號
	Buyer              string `json:"Buyer"`              //會員帳號
	TerminalId         string `json:"TerminalId"`         //會員代碼
	WithdrawAmount     int64  `json:"Amount"`             //提領金額
	WithdrawFee        int64  `json:"Fee"`                //手續費
	BankAccount        string `json:"BankAccount"`        //銀行帳號
	BankName           string `json:"BankName"`           //銀行
	WithdrawType       string `json:"WithdrawType"`       //匯款方式
	WithdrawStatus     string `json:"WithdrawStatus"`     //提領狀態
	WithdrawStatusText string `json:"WithdrawStatusText"` //提領狀態
}

type SearchProductResponse struct {
	CreateDate      string `json:"CreateDate"`      //上架日期          開立日期
	UpdateDate      string `json:"UpdateDate"`      //修改日期          －－－        －－－
	UserAccount     string `json:"UserAccount"`     //賣家帳號                       買家帳號
	StoreName       string `json:"StoreName"`       //收銀機名稱                     買家暱稱
	ProductId       string `json:"ProductId"`       //商品編號          帳單編號      訂購單編號
	ProductName     string `json:"ProductName"`     //商品名稱                       訂購內容
	Amount          int64  `json:"Amount"`          //單價             帳單金額       單價
	Quantity        int64  `json:"Quantity"`        //可賣數量          －－－        數量
	ShipFee         int64  `json:"ShipFee"`         //                              運費
	Total           int64  `json:"Total"`           //                              總金額
	ProductSpec     string `json:"ProductSpec"`     //規格
	ShipMode        string `json:"ShipMode"`        //出貨方式
	PayWayMode      string `json:"PayWayMode"`      //付款方式
	FreeShipping    string `json:"FreeShipping"`    //免運費
	SetFreeShip     string `json:"SetFreeShip"`     //免運費設定
	ShipMerge       string `json:"ShipMerge"`       //合併運費
	GoogleFormLink  string `json:"GoogleFormLink"`  //Google表單連結
	ProductLink     string `json:"ProductLink"`     //商品頁網址
	BuyerName       string `json:"BuyerName"`       //訂購人
	ReceiverName    string `json:"ReceiverName"`    //收件人
	ReceiverAddress string `json:"ReceiverAddress"` //收件地址
	Url             string `json:"Url"`
}

type CustomerContactResponse struct {
	ID           string             `json:"id"`
	RelatedId    string             `json:"relatedId"`
	QuestionId   string             `json:"questionId"`
	OrderID      string             `json:"orderId"`
	Type         string             `json:"type"`
	Content      string             `json:"content"`
	CreateTime   string             `json:"createTime"`
	ReplyContent string             `json:"replyContent"`
	ReplyTime    string             `json:"replyTime"`
	Member       CustomerMemberInfo `json:"member"`
	Memos        []CustomerMemoInfo `json:"memos"`
}

type CustomerMemberInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Account string `json:"account"`
}

type OrderWithProducts struct {
	CreateDate        string               `json:"createDate"`
	OrderID           string               `json:"orderId"`
	Status            string               `json:"status"`
	TransFee          string               `json:"transFee"`
	InfoFee           string               `json:"infoFee"`
	PayFee            string               `json:"payFee"`
	PlatShipFee       string               `json:"platShipFee"`
	TotoalAmount      string               `json:"totoalAmount"`
	ProductTotalPrice string               `json:"productTotalPrice"`
	PlatFee           string               `json:"platFee"`
	ShipFee           string               `json:"shipFee"`
	ShipMerge         bool                 `json:"shipMerge"`
	Products          []OrderProduct       `json:"products"`
	Messages          OrderMessage         `json:"messages"`
	Rufund            []OrderRefund        `json:"refund"`
	ReturnProduct     []OrderReturnProduct `json:"returnProducts"`
}

type OrderProduct struct {
	ShipFee  string `json:"shipFee"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type OrderMessage struct {
	Count        int           `json:"messageCount"`
	OrderMessage []MessageData `json:"messages"`
}

type MessageData struct {
	OrderID    string `json:"orderId"`
	CreateTime string `json:"createTime"`
	Role       string `json:"role"`
	Message    string `json:"message"`
}

type OrderReturnProduct struct {
	Status     string `json:"status"`
	RefundTime string `json:"refundTime"`
	Quantity   int    `json:"qty"`
	ProductID  string `json:"type"`
	Name       string `json:"name"`
	ID         string `json:"id"`
	Price      string `json:"price"`
}

type OrderRefund struct {
	Status     string `json:"status"`
	RefundTime string `json:"refundTime"`
	ID         string `json:"id"`
	Amount     string `json:"price"`
}

type ReplyCustomer struct {
	Status string `json:"status"`
}

type CustomerMemo struct {
	Status string `json:"status"`
}

type CustomerMemoInfo struct {
	ID         string `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	Staff      string `json:"staff,omitempty"`
	CreateTime string `json:"createdTime,omitempty"`
}

type SearchMemberResponse struct {
	Count   int64             `json:"Count"`
	Members []MembersResponse `json:"Orders"`
}

type MembersResponse struct {
	Account         string   `json:"Account"`         //會員帳號
	Nickname        string   `json:"Nickname"`        //會員暱稱
	TerminalId      string   `json:"TerminalId"`      //會員代碼
	Uid             string   `json:"Uid"`             //會員id
	RegisterTime    string   `json:"RegisterTime"`    //加入會員日期
	Email           string   `json:"Email"`           //Email
	Upgrade         string   `json:"Upgrade"`         //付費方案
	UpgradeCycle    string   `json:"UpgradeCycle"`    //付費週期
	Outstanding     int64    `json:"Outstanding"`     //未付訂單金額
	Store           int64    `json:"Store"`           //賣場數
	Manage          int64    `json:"Manage"`          //管理帳號數
	Status          string   `json:"Status"`          //帳號狀態
	StoreName       []string `json:"StoreName"`       //賣場
	Balance         int64    `json:"Balance"`         //餘額
	DetainBalance   int64    `json:"DetainBalance"`   //扣留餘額
	WithholdBalance int64    `json:"WithholdBalance"` //保留餘額
	RetainBalance   int64    `json:"RetainBalance"`   //待撥付餘額
}

type StoreManager struct {
	Managers []RankAccount `json:"Managers"`
}
type RankAccount struct {
	MainAccount string `json:"MainAccount"`
	Account     string `json:"Account"`
	NickName    string `json:"NickName"`
	Status      string `josn:"Status"`
	StoreName   string `json:"StoreName"`
	Email       string `json:"Email"`
	CreateTime  string `json:"CreateTime"`
	DeleteTime  string `json:"DeleteTime"`
}

//收費紀錄
type MemberMainResponse struct {
	Master  MemberMaster `json:"Master"`
	Manager StoreManager `json:"Manager"`
}
type MemberMaster struct {
	Stores []MemberStore `json:"Stores"`
}
type MemberStore struct {
	StoreName    string        `json:"StoreName"`
	StoreStatus  string        `json:"StoreStatus"`
	ProductCount int64         `json:"ProductCount"`
	InstantCount int64         `json:"InstantCount"`
	Created      string        `json:"Created"`
	Updated      string        `json:"Updated"`
	Deleted      []RankAccount `json:"Deleted"`
	Operated     []RankAccount `json:"Operated"`
}

//收費紀錄
//賣場-主帳號 列表
//運作中管理帳號 列表
//已刪除管理帳號 列表
//賣場-管理帳號 列表
