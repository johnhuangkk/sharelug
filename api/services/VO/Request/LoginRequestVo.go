package Request

type PaySendOtpParams struct {
	Phone string  `form:"Phone" validate:"required"`
}

type SendOtpParams struct {
	Phone string `form:"Phone" validate:"required"`
	Login bool   `form:"Login"`
}

type ValidateOtpParams struct {
	Phone string `form:"Phone" validate:"required"`
	Code  string `form:"Code" validate:"required"`
}

type ExchangeStoreParams struct {
	StoreId string `form:"StoreId" json:"StoreId" validate:"required"`
}