package Response

type BillListResponse struct {
	Tabs     BillTabs       `json:"Tabs"`
	Count    int64          `json:"Count"`
	BillList []BillResponse `json:"BillList"`
}

type BuyerBillListResponse struct {
	Count    int64          `json:"Count"`
	BillList []BillResponse `json:"BillList"`
}


type BillTabs struct {
	Valid   int64 `json:"Valid"`
	Cancel  int64 `json:"Cancel"`
	Overdue int64 `json:"Overdue"`
}

type BillResponse struct {
	BillId           string     `json:"BillId"`
	ProductImage     string     `json:"ProductImage"`
	ProductName      string     `json:"ProductName"`
	ProductLink      string     `json:"ProductLink"`
	ProductSpec      string     `json:"ProductSpec"`
	ProductPrice     int64      `json:"ProductPrice"`
	ProductQty       int64      `json:"ProductQty"`
	BuyerName        string     `json:"BuyerName"`
	ReceiverName     string     `json:"ReceiverName"`
	ReceiverAddress  string     `json:"Address"`
	ShipType         string     `json:"ShipType"`
	ShipFee          int64      `json:"ShipFee"`
	PayWayType       string     `json:"PayWayType"`
	TotalAmount      int64      `json:"TotalAmount"`
	PlatformShipFee  int64      `json:"PlatformShipFee"`
	PlatformTransFee int64      `json:"PlatformTransFee"`
	PlatformInfoFee  int64      `json:"PlatformInfoFee"`
	PlatformPayFee   int64      `json:"PlatformPayFee"`
	CaptureAmount    int64		`json:"CaptureAmount"`	
	TinyUrl          string     `json:"TinyUrl"`
	Qrcode           string     `json:"Qrcode"`
	Expire           string     `json:"Expire"`
	IsExtension      bool       `json:"IsExtension"`
	BillStatus       string     `json:"BillStatus"`
	PayWayStatus     string		`json:"PayWayStatus"`
	MemberInfo       MemberInfo `json:"MemberInfo"`
}

type BillConfirmResponse struct {
	OrderId string `json:"OrderId"`
}
