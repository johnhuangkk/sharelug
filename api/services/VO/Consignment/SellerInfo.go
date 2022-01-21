package Consignment

// 託運單賣家資訊
type SellerInfo struct {
	OrderId []string `form:"orderId" json:"orderId" valid:"required"`
}
