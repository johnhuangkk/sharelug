package Response

type PaymentResponse struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
	RtnHtml string `json:"RtnHtml"`
}

type PayResponse struct {
	OrderId string `json:"OrderId"`
	RtnURL  string `json:"RtnURL"`
}

type OrderResponse struct {
	OrderId           string              `json:"OrderId"`
	StoreName         string              `json:"StoreName"`
	StoreId           string              `json:"StoreId"`
	OrderTime         string              `json:"OrderTime"`
	Buyer             Buyer               `json:"Buyer"`
	BuyerMasker       Buyer               `json:"BuyerMasker"`
	PayWayTime        string              `json:"PayWayTime"`
	CaptureTime       string              `json:"CaptureTime"`
	TotalShipFee      int64               `json:"TotalShipFee"`
	Coupon            int64               `json:"Coupon"`
	CouponNumber      string              `json:"CouponNumber"`
	TotalAmount       int64               `json:"TotalAmount"`
	RefundAmount      int64               `json:"RefundAmount"`
	SubTotal          int64               `json:"SubTotal"`
	ShipType          string              `json:"ShipType"`
	ShipCompany       string              `json:"ShipCompany"`
	PlatformShipFee   int64               `json:"PlatformShipFee"`
	PlatformTransFee  int64               `json:"PlatformTransFee"`
	PlatformInfoFee   int64               `json:"PlatformInfoFee"`
	PlatformPayFee    int64               `json:"PlatformPayFee"`
	Income            int64               `json:"Income"`
	OrderStatusType   string              `json:"OrderStatusType"`
	OrderStatusText   string              `json:"OrderStatusText"`
	ShipStatusType    string              `json:"ShipStatusType"`
	ShipStatusText    string              `json:"ShipStatusText"`
	IsReturn          bool                `json:"IsReturn"`
	IsRefund          bool                `json:"IsRefund"`
	CaptureStatusType string              `json:"CaptureStatusType"`
	CaptureStatusText string              `json:"CaptureStatusText"`
	CaptureApplyType  string              `json:"CaptureApplyType"`
	CaptureApplyText  string              `json:"CaptureApplyText"`
	ShipExpire        string              `json:"ShipExpire"`
	ShipNumber        string              `json:"ShipNumber"`
	ShipTime          string              `json:"ShipTime"`
	InvoiceNumber     string              `json:"InvoiceNumber"`
	RefundStatus      string              `json:"RefundStatus"`
	OrderMemo         string              `json:"OrderMemo"`
	BuyerNotes        string              `json:"BuyerNotes"`
	FormUrl           string              `json:"FormUrl"`
	Detail            OrderDetailResponse `json:"Detail"`
	Payment           interface{}         `json:"Payment"`
	Shipping          Shipping            `json:"Shipping"`
	PayType           string              `json:"PayType"`
}

type Buyer struct {
	BuyerName  string `json:"BuyerName"`
	BuyerPhone string `json:"BuyerPhone"`
	BuyerUid   string `json:"BuyerId"`
}

type Shipping struct {
	Type           string   `json:"Type"`
	Text           string   `json:"Text"`
	Receiver       Receiver `json:"Receiver"`
	ReceiverMasker Receiver `json:"ReceiverMasker"`
}

type Receiver struct {
	ReceiverName    string `json:"ReceiverName"`
	ReceiverPhone   string `json:"ReceiverPhone"`
	ReceiverAlias   string `json:"ReceiverAlias"`
	ReceiverAddress string `json:"ReceiverAddress"`
}

type OrderDetailResponse struct {
	Merge      []OrderDetails `json:"CanNotMergeList"`
	MergeFee   int            `json:"CanNotMergeFee"`
	General    []OrderDetails `json:"GeneralList"`
	GeneralFee int            `json:"GeneralFee"`
	Free       []OrderDetails `json:"ShipFreeList"`
	FreeFee    int            `json:"ShipFreeFee"`
}

type OrderDetails struct {
	ProductId       string `json:"ProductId"`
	ProductSpecName string `json:"ProductSpec"`
	ProductSpecId   string `json:"ProductSpecId"`
	ProductName     string `json:"ProductName"`
	ProductQuantity int64  `json:"Quantity"`
	ProductPrice    int64  `json:"ProductPrice"`
	ShipMerge       int64  `json:"ShipMerge"`
	ShipFee         int64  `json:"ShipFee"`
	ReturnQty       int64  `json:"ReturnQty"`
}

type B2COrderResponse struct {
	OrderId       string        `json:"OrderId"`     //訂單編號
	ProductName   string        `json:"ProductName"` //方案名稱
	OrderTime     string        `json:"OrderTime"`   //訂購時間
	ExpireTime    string        `json:"ExpireTime"`  //到期時間
	UpgradeList   []UpgradeList `json:"UpgradeList"`
	OrderList     []OrderList   `json:"OrderList"`
	UpgradeSum    int64         `json:"UpgradeSum"`
	OrderSum      int64         `json:"OrderSum"`
	PriceTotal    int64         `json:"PriceTotal"`    //結帳金額
	Payment       string        `json:"Payment"`       //付款方式
	BankCode      string        `json:"BankCode"`      //銀行代碼
	BankAccount   string        `json:"BankAccount"`   //銀行帳號
	AtmExpire     string        `json:"AtmExpire"`     //ATM到期時間
	ProductDetail ProductDetail `json:"ProductDetail"` //方案內容
}

//商品名稱、單價、類別
type B2cOrderDetail struct {
	ProductId     string `json:"ProductId"`
	ProductName   string `json:"ProductName"`
	ProductDetail string `json:"ProductDetail"`
	ProductAmount int64  `json:"ProductAmount"`
	ProductType   string `json:"ProductType"`
	BillingTime   string `json:"BillingTime"`
}

type ProductDetail struct {
	ProductAmt int64 `json:"ProductAmt"`
	Store      int64 `json:"Store"`
	Manager    int64 `json:"Manager"`
}
