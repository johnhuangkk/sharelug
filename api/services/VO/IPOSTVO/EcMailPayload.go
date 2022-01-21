package IPOSTVO

/**
賣家出貨時需填寫之資料
*/
type SellerShipOrder struct {
	OrderId []string `form:"orderId" json:"orderId" valid:"required"`
}
