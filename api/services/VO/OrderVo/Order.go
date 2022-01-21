package OrderVo


type CalculatePlatFormFeeVo struct {
	ShipType   string
	PayWayType string
	Amount     float64
}

type CreditPaymentVo struct {
	OrderId     string
	TotalAmount float64
	OrderType   string
	BuyerId     string
	AuditStatus string
	SellerId    string
}

type CalculatePlatFormFeeResponse struct {
	PlatformTransFee float64
	PlatformShipFee  float64
	PlatformPayFee   float64
	PlatformInfoFee  float64
	CaptureAmount    float64
}
