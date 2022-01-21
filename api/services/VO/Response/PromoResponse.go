package Response

type PromoCreate struct {
	Id        string `json:"Id"`
	Name      string `json:"Name"`
	Quantity  int64  `json:"Quantity"`
	Amount    int64  `json:"Amount"`
	Status    string `json:"Status"`
	StartTime string `json:"StartTime"`
	StopTime  string `json:"StopTime"`
	EndTime   string `json:"EndTime"`
}

type PromoDetail struct {
	Id         int64  `json:"Id"`
	Name       string `json:"Name"`
	Quantity   int64  `json:"Quantity"`
	Remain     int64  `json:"Remain"`
	Picked     int64  `json:"Picked"`
	Used       int64  `json:"Used"`
	UnUsed     int64  `json:"UnUsed"`
	Amount     int64  `json:"Amount"`
	Status     string `json:"Status"`
	StatusText string `json:"StatusText"`
	StartTime  string `json:"StartTime"`
	StopTime   string `json:"StopTime"`
	EndTime    string `json:"EndTime"`
}

type PromoList struct {
	PromoEnable bool          `json:"Enable"`
	Promos      []PromoDetail `json:"Promos"`
}

type CouponUnuse struct {
	Id          int64  `json:"Id"`
	PromotionId int64  `json:"PromotionId"`
	Code        string `json:"Code"`
	IsCopy      bool   `json:"IsCopy"`
	StatusText  string `json:"StatusText"`
}

type CouponUsed struct {
	Code           string  `json:"Code"`
	StatusText     string  `json:"StatusText"`
	OrderId        string  `json:"OrderId"`
	BuyerPhone     string  `json:"BuyerPhone"`
	UseDate        string  `json:"UseDate"`
	Amount         float64 `json:"Amount"`
	DiscountAmount float64 `json:"DiscountAmount"`
}

type CouponUsedList struct {
	Counts  int64        `json:"Counts"`
	Coupons []CouponUsed `json:"Coupons"`
}

type CouponUnuseList struct {
	Counts  int64         `json:"Counts"`
	Coupons []CouponUnuse `json:"Coupons"`
}
