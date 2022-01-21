package Request

type PromoCreate struct {
	Name     string `json:"Name"` //二十個字
	Quantity int64  `json:"Quantity"`
	Amount   int64  `json:"Amount"`
	EndTime  string `json:"EndTime"`
}

type PromoEnable struct {
	Enable bool `json:"Enable"`
}

type PromoStop struct {
	PromoId string `json:"PromoId"`
}

type PromoCouponTake struct {
	PromoId  string `json:"PromoId"`
	Quantity string `json:"Quantity"`
}

type PromoCouponListGet struct {
	Page  int `json:"Page"`
	Limit int `json:"Limit"`
}
