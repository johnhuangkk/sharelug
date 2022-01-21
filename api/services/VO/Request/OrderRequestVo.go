package Request

type PayParams struct {
	BuyerPhone     string `form:"BuyerPhone" validate:"required"`
	BuyerName      string `form:"BuyerName"`
	ReceiverId     string `form:"ReceiverId"`
	ReceiverName   string `form:"ReceiverName"`
	ReceiverPhone  string `form:"ReceiverPhone"`
	Address        string `form:"Address"`
	CardId         string `form:"CreditId" json:"CreditId"`
	CardNumber     string `form:"CreditNumber" json:"CreditNumber"`
	CardExpiration string `form:"CreditExpiration" json:"CreditExpiration"`
	CardSecurity   string `form:"CreditSecurity" json:"CreditSecurity"`
	PayWay         string `form:"PayWay" json:"PayWay" validate:"required"`
	BuyerNotes     string `form:"BuyerNotes" json:"BuyerNotes"`
	PayType        string `form:"PayType" json:"PayType"`
}

type OrderReadParams struct {
	OrderId []string `form:"orderId" validate:"required"`
}

type OrderMemoParams struct {
	OrderId   string `form:"OrderId" validate:"required"`
	OrderMemo string `form:"OrderMemo"`
}

type SetShipNumberParams struct {
	ShipNumberList []ShipNumberList `form:"ShipNumberList" json:"ShipNumberList"`
}

type ExportOrderShippingParams struct {
	ShipType  string   `form:"ShipType" json:"ShipType"`
	Orders    []string `form"Orders" json:"Orders"`
	OrderBy   string   `form:"OrderBy" json:"OrderBy"`
	SelectAll bool     `form:"SelectAll" json:"SelectAll"`
}

type ShipNumberList struct {
	OrderId  string `form:"OrderId" json:"OrderId" validate:"required"`
	ShipText string `form:"ShipText" json:"ShipText"`
	Number   string `form:"Number" json:"Number" validate:"required"`
}

type SetPaymentParams struct {
	OrderId string `form:"OrderId" json:"OrderId" validate:"required"`
}

type Credit3dCheckParams struct {
	MerchantID    string `form:"MerchantID" json:"MerchantID"`
	TerminalID    string `form:"TerminalID" json:"TerminalID"`
	OrderID       string `form:"OrderID" json:"OrderID"`
	PAN           string `form:"PAN" json:"PAN"`
	TransCode     string `form:"TransCode" json:"TransCode"`
	TransMode     string `form:"TransMode" json:"TransMode"`
	TransDate     string `form:"TransDate" json:"TransDate"`
	TransTime     string `form:"TransTime" json:"TransTime"`
	TransAmt      string `form:"TransAmt" json:"TransAmt"`
	ApproveCode   string `form:"ApproveCode" json:"ApproveCode"`
	ResponseCode  string `form:"ResponseCode" json:"ResponseCode"`
	ResponseMsg   string `form:"ResponseMsg" json:"ResponseMsg"`
	InstallType   string `form:"InstallType" json:"InstallType"`
	Install       string `form:"Install" json:"Install"`
	FirstAmt      string `form:"FirstAmt" json:"FirstAmt"`
	EachAmt       string `form:"EachAmt" json:"EachAmt"`
	Fee           string `form:"Fee" json:"Fee"`
	RedeemType    string `form:"RedeemType" json:"RedeemType"`
	RedeemUsed    string `form:"RedeemUsed" json:"RedeemUsed"`
	RedeemBalance string `form:"RedeemBalance" json:"RedeemBalance"`
	CreditAmt     string `form:"CreditAmt" json:"CreditAmt"`
	RiskMark      string `form:"RiskMark" json:"RiskMark"`
	FOREIGN       string `form:"FOREIGN" json:"FOREIGN"`
	SECURE_STATUS string `form:"SECURE_STATUS" json:"SECURE_STATUS"`
	PrivateData   string `form:"PrivateData" json:"PrivateData"`
	KEY           string `form:"KEY" json:"KEY"`
	Signature     string `form:"Signature" json:"Signature"`
}
