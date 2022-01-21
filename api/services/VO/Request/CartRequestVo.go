package Request

type BillParams struct {
	ProductSpecId 	string 	`form:"ProductSpecId" validate:"required"`	//商品編號
	Quantity  		int 	`form:"Quantity" validate:"required"`		//商品數量
	Shipping		string	`form:"shipping"`							//運送方式
}

type AddCartParams struct {
	ProductSpecId 	string 	`form:"ProductSpecId" validate:"required"`	//商品編號
	Quantity  		int 	`form:"Quantity" validate:"required"`		//商品數量
	Shipping		string	`form:"shipping"`							//運送方式
}

type DeleteCartParams struct {
	ProductSpecId string	`form:"ProductSpecId" validate:"required"`
}

type ChangeShippingParams struct {
	ShipType string	`form:"ShipType" validate:"required"`
}

type ChangeQuantityParams struct {
	Type 			string 	`form:"Type" validate:"required"`
	ProductSpecId	string	`form:"ProductSpecId" validate:"required"`
}

type GetAddressParams struct {
	ShipType string `form:"ShipType" validate:"required"`
}

type CheckCouponNumberParams struct {
	CouponNumber string `form:"CouponNumber" validate:"required"`
}