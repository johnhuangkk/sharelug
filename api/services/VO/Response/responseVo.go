package Response

type OrderSearchResponse struct {
	Tabs        OrderTabsResponse   `json:"Tabs"`
	OrderList   []OrderListResponse `json:"OrderList"`
	OrderCount  int64               `json:"OrderCount"`
	UnreadCount int64               `json:"UnreadCount"`
}

type OrderTabsResponse struct {
	All         int64 `json:"All"`
	OrderWait   int64 `json:"OrderWait"`
	ShipWait    int64 `json:"ShipWait"`
	Expire      int64 `json:"Expire"`
	OrderCancel int64 `json:"OrderCancel"`
}

type ShipSearchResponse struct {
	Tabs        ShipTabsResponse    `json:"Tabs"`
	OrderList   []OrderListResponse `json:"OrderList"`
	OrderCount  int64               `json:"OrderCount"`
	UnreadCount int64               `json:"UnreadCount"`
}

type ShipTabsResponse struct {
	ShipWait    int64 `json:"ShipWait"`
	ShipOverdue int64 `json:"ShipOverdue"`
	Shipment    int64 `json:"Shipment"`
	ShipSucc    int64 `json:"ShipSucc"`
}

type BuyerOrderSearchResponse struct {
	Tabs        BuyerOrderTabsResponse `json:"Tabs"`
	OrderList   []OrderListResponse    `json:"OrderList"`
	OrderCount  int64                  `json:"OrderCount"`
	UnreadCount int64                  `json:"UnreadCount"`
}

type BuyerOrderTabsResponse struct {
	All         int64 `json:"All"`
	OrderWait   int64 `json:"OrderWait"`
	Shipment    int64 `json:"Shipment"`
	OrderCancel int64 `json:"OrderCancel"`
	Refund      int64 `json:"Refund"`
}

type OrderListResponse struct {
	OrderId           string        `json:"orderId"`
	StoreId           string        `json:"storeId"`
	BuyerId           string        `json:"buyerId"`
	OrderStatusType   string        `json:"orderStatusType"`
	OrderStatusText   string        `json:"orderStatusText"`
	ShipStatusType    string        `json:"shipStatusType"`
	ShipStatusText    string        `json:"shipStatusText"`
	RefundStatusType  string        `json:"refundStatusType"`
	RefundStatusText  string        `json:"refundStatusText"`
	CaptureStatusType string        `json:"captureStatusType"`
	CaptureStatusText string        `json:"captureStatusText"`
	ShipNumber        string        `json:"shipNumber"`
	ShipTime          string        `json:"shipTime"`
	ShipExpire        string        `json:"shipExpire"`
	ShipCompany       string        `json:"shipCompany"`
	CreateTime        string        `json:"createTime"`
	CaptureTime       string        `json:"captureTime"`
	PayWayType        string        `json:"payWayType"`
	PayWayTime        string        `json:"payWayTime"`
	PayWayText        string        `json:"payWayText"`
	InvoiceNumber     string        `json:"invoiceNumber"`
	Detail            []OrderDetail `json:"detail"`
	BuyerName         string        `json:"buyerName"`
	StoreName         string        `json:"StoreName"`
	Ship              string        `json:"ship"`
	ShipType          string        `json:"shipType"`
	Price             int64         `json:"Price"`
	Unread            bool          `json:"unread"`
	FormUrl           string        `json:"FormUrl"`
	BuyerNotes        string        `json:"BuyerNotes"`
}

type OrderDetail struct {
	ProductName string `json:"productName"`
}

type NotificationResponse struct {
	Tabs          NotifyTabsResponse `json:"Tabs"`
	NotifyMessage []NotifyResponse   `json:"NotifyMessage"`
	NotifyCount   int64              `json:"NotificationCount"`
}

type NotifyTabsResponse struct {
	System   int64 `json:"System"`
	Platform int64 `json:"Platform"`
}

type NotifyResponse struct {
	MessageId  int64  `json:"MessageId"`
	Message    string `json:"Message"`
	MsgType    string `json:"MsgType"`
	OrderId    string `json:"OrderId"`
	CreateTime string `json:"CreateTime"`
	Unread     bool   `json:"Read"`
}

type UpgradePlanResponse struct {
	StoreName           string            `json:"StoreName"`
	StorePicture        string            `json:"StorePicture"`
	StoreManagerLimit   int64             `json:"StoreLimit"`
	StoreManagerCurrent int64             `json:"StoreCurrent"`
	StoreUpgradeText    string            `json:"StoreUpgradeText"`
	CurrentPlan         UpgradePlanList   `json:"CurrentPlan"`
	UpgradePlanList     []UpgradePlanList `json:"UpgradePlanList"`
}

type UpgradePlanList struct {
	ProductId   string   `json:"ProductId"`
	ProductName string   `json:"ProductName"`
	Description []string `json:"Description"`
	Note        string   `json:"Note"`
	Amount      int64    `json:"Amount"`
	IsPay       bool     `json:"IsPay"`
}

type GetB2cPayResponse struct {
	UpgradeList []UpgradeList `json:"UpgradeList"`
	OrderList   []OrderList   `json:"OrderList"`
	UpgradeSum  int64         `json:"UpgradeSum"`
	OrderSum    int64         `json:"OrderSum"`
	PriceTotal  int64         `json:"PriceTotal"`
	OrderId     string        `json:"OrderId"`
}

type UpgradeList struct {
	UpgradeText  string `json:"UpgradeText"`
	UpgradePrice int64  `json:"UpgradePrice"`
	SignType     bool   `json:"SignType"`
}

type OrderList struct {
	OrderText  string `json:"OrderText"`
	OrderPrice int64  `json:"OrderPrice"`
}

type B2cPayResponse struct {
	OrderId string `json:"OrderId"`
	RtnURL  string `json:"RtnURL"`
}

type ManagerListResponse struct {
	StoreLimit   int64         `json:"StoreLimit"`
	StoreCurrent int64         `json:"StoreCurrent"`
	ManagerList  []ManagerList `json:"ManagerList"`
}

type ManagerList struct {
	ManagerId         int    `json:"ManagerId"`
	ManagerPicture    string `json:"ManagerPicture"`
	ManagerName       string `json:"ManagerName"`
	ManagerEmail      string `json:"ManagerEmail"`
	ManagerStatus     string `json:"ManagerStatus"`
	ManagerStatusText string `json:"ManagerStatusText"`
	ManagerStartTime  string `json:"ManagerStartTime"`
}

type MemberEmailVerifyResponse struct {
	VerifyType      string `json:"VerifyType"`
	VerifyStoreName string `json:"VerifyStoreName"`
	VerifyStatus    string `json:"VerifyStatus"`
}
type MemberCompanyVerifyResponse struct {
	Mphone         string `json:"Mphone"`
	Representative string `json:"Representative"`
	IdentityName   string `json:"IdentityName"`
	CompanyName    string `json:"CompanyName"`
	CompanyAddr    string `json:"CompanyAddr"`
	Category       string `json:"Category"`
	VerifyStatus   string `json:"VerifyStatus"`
}

type LoginInfoResponse struct {
	Member MemberInfo `json:"Member"`
	Store  StoreInfo  `json:"Store"`
}

type CountUpgradeOrderResponse struct {
	UpgradeCount   int64  `json:"UpgradeCount"`
	BillingOrderId string `json:"BillingOrderId"`
}

type UnpaidUpgradeOrderResponse struct {
	OrderId       string `json:"OrderId"`
	UserId        string `json:"UserId""`
	StoreId       string `json:"StoreId"`
	ProductId     string `json:"ProductId"`
	ProductName   string `json:"ProductName"`
	ProductDetail string `json:"ProductDetail"`
	BillingTime   string `json:"BillingTime"`
	Amount        int64  `json:"Amount"`
	Payment       string `json:"Payment"`
	Status        string `json:"Status"`
	CreateTime    string `json:"CreateTime"`
}

type NoticeResponse struct {
	Message int64 `json:"Message"`
}

type CarrierResponse struct {
	InvoiceType string               `json:"InvoiceType"`
	CompanyBan  string               `json:"CompanyBan"`
	CompanyName string               `json:"CompanyName"`
	DonateBan   string               `json:"DonateBan"`
	CarrierType string               `json:"CarrierType"`
	CarrierId   string               `json:"CarrierId"`
	DonateUnit  []DonateListResponse `json:"DonateUnit"`
	Choose      string               `json:"Choose"`
}

type DonateListResponse struct {
	DonateName string `json:"DonateName"`
	DonateCode string `json:"DonateCode"`
}

type BindCarrierResponse struct {
	CardBan   string `json:"card_ban"`
	CardNo1   string `json:"card_no1"`
	CardNo2   string `json:"card_no2"`
	CardType  string `json:"card_type"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
	Action    string `json:"action"`
}

type VerifyDonateCodeResponse struct {
	Version string `json:"v"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	IsExist string `json:"isExist"`
	TxID    string `json:"TxID"`
}

type VerifyCompanyBan struct {
	Version         string `json:"v"`
	Code            int    `json:"code"`
	Msg             string `json:"msg"`
	HashSerial      string `json:"hashSerial"`
	BanUnitTpStatus string `json:"banUnitTpStatus"`
}

type ShipReportPdfResponse struct {
	Id              int64
	OrderId         string
	ReceiverName    string
	ReceiverPhone   string
	ReceiverCode    string
	ReceiverCity    string
	ReceiverArea    string
	ReceiverAddress string
}

type DeliveryReportResponse struct {
	StoreName       string           `json:"StoreName"`
	OrderId         string           `json:"OrderId"`
	OrderDate       string           `json:"OrderDate"`
	ReceiverName    string           `json:"ReceiverName"`
	ReceiverPhone   string           `json:"ReceiverPhone"`
	ReceiverCode    string           `json:"ReceiverCode"`
	ReceiverCity    string           `json:"ReceiverCity"`
	ReceiverArea    string           `json:"ReceiverArea"`
	ReceiverAddress string           `json:"ReceiverAddress"`
	OrderMemo       string           `json:"OrderMemo"`
	BuyerNotes      string           `json:"BuyerNotes"`
	Details         []DeliveryDetail `json:"Details"`
	ShipFee         int64            `json:"ShipFee"`
	TotalAmount     int64            `json:"TotalAmount"`
}

type DeliveryDetail struct {
	Id          int64  `json:"Id"`          //項目
	ProductName string `json:"ProductName"` //品名
	Quantity    int64  `json:"Quantity"`    //數量
	Price       int64  `json:"Price"`       //單價
}

type F2fReportPdfResponse struct {
	OrderId       string       `json:"OrderId"`
	ReceiverName  string       `json:"ReceiverName"`
	ReceiverPhone string       `json:"ReceiverPhone"`
	Products      []F2fProduct `json:"Products"`
}

type F2fProduct struct {
	ProductName     string `json:"ProductName"`
	ProductQuantity int64  `json:"ProductQuantity"`
}

type BatchOrderShippingResponse struct {
	BatchId     string        `json:"BatchId"`
	BatchOrders []BatchOrders `json:"BatchOrders"`
}

type BatchOrders struct {
	Id      int64  `json:"Id"`
	OrderId string `json:"OrderId"`
	Trader  string `json:"Trader"`
	Number  string `json:"Number"`
}

type CityWithArea struct {
	CityCode string                `json:"CityCode"`
	CityName string                `json:"CityName"`
	Area     []AreaFlagWithZipCode `json:"Area"`
}

type AreaWithZipCode struct {
	ZipCode string `json:"ZipCode"`
	Name    string `json:"Name"`
}

type CityWithCode struct {
	CityCode string `json:"CityCode"`
	CityName string `json:"CityName"`
}

type AreaFlagWithZipCode struct {
	ZipCode string `json:"ZipCode"`
	Name    string `json:"Name"`
	Enable  bool   `json:"Enable"`
}

type CompanyVerifyPendingList struct {
	Counts    int64               `json:"Counts"`
	Companies []CompanyVerifyInfo `json:"Companies"`
}
type CompanyVerifyInfo struct {
	UserId           string `json:"UserId"`
	UserPhone        string `json:"UserPhone"`
	Representative   string `json:"Representative"` //統一編號
	RepresentativeId string `json:"RepresentativeId"`
	CompanyName      string `json:"CompanyName"`
}

type MemberSpecialStore struct {
	Account            string  `json:"Account"`
	CompanyName        string  `json:"CompanyName"`
	RepresentativeId   string  `json:"RepresentativeId"`
	Representative     string  `json:"Representative"`
	RepresentFirst     string  `json:"RepresentFirst"`
	RepresentLast      string  `json:"RepresentLast"`
	MemberPhone        string  `json:"MemberPhone"`
	MemberName         string  `json:"MemberName"`
	IdentityId         string  `json:"IdentityId"`
	Capital            float64 `json:"Capital"`
	Establish          string  `json:"Establish"`
	ZipCode            string  `json:"ZipCode"`
	Addr               string  `json:"Addr"`
	AddrEn             string  `json:"AddrEn"`
	Contact            string  `json:"Contact"`
	ContactPhone       string  `json:"ContactPhone"`
	SpecialStoreName   string  `json:"SpecialStoreName"`
	SpecialStoreNameEn string  `json:"SpecialStoreNameEn"`
	MerchantId         string  `json:"MerchantId"`
	Terminal3D         string  `json:"Terminal3D"`
	TerminalId         string  `json:"TerminalId"`
	JobCode            string  `json:"JobCode"`
	MccCode            string  `json:"MccCode"`
	CityCode           string  `json:"CityCode"`
	Created            string  `json:"Created"`
}

type KgiSpecialStoreFile struct {
	Name    string `json:"Name"`
	Created string `json:"Created"`
}
