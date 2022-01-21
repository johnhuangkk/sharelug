package Request

type SalesReportRequest struct {
	Tab       string `json:"Tab"`
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
}

type SettingStoreRequest struct {
	StoreStatus  int    `json:"StoreStatus"`
	StoreName    string `json:"StoreName"`
	StorePicture string `json:"StorePicture"`
}

type SettingFreeShipRequest struct {
	FreeShipKey string `json:"FreeShipKey"`
	FreeShip    int64  `json:"FreeShip"`
}

type SettingUserRequest struct {
	UserEmail   string `form:"UserEmail" json:"UserEmail"`
	UserName    string `form:"UserName" json:"UserName"`
	UserPicture string `form:"UserPicture" json:"UserPicture"`
}

type SetRealTimesRequest struct {
	ProductId string `json:"ProductId"`
}

type BalanceRequest struct {
	Tab       string `json:"Tab"`
	StartTime string `json:"StartTime"`
	EndTime   string `json:"EndTime"`
	Limit     int64  `json:"Limit"`
	Page      int64  `json:"Page"`
}

type MyAccountRequest struct {
	Limit int64 `json:"Limit"`
	Page  int64 `json:"Page"`
}

type RetainRequest struct {
	Limit int64 `json:"Limit"`
	Page  int64 `json:"Page"`
}

type StoreSocialMediaRequest struct {
	Type string `form:"Type" json:"Type"`
	Link string `form:"Link" json:"Link"`
	Show int    `form:"Show" json:"Show"`
}

type SelfDeliveryAreaRequest struct {
	StoreId string         `json:"StoreId"`
	Enable  bool           `json:"Enable"`
	Section []DeliveryArea `json:"Section"`
}

type DeliveryArea struct {
	CityCode string   `json:"CityCode"`
	AreaList []string `json:"AreaList"`
}

type Area struct {
	ZipCode string `json:"ZipCode"`
}

type SelfDeliveryFeeRequest struct {
	StoreId                 string `json:"storeId"`
	SelfDeliveryFreeShipKey string `json:"SelfDeliveryFreeShipKey"`
	SelfDeliveryFree        int64  `json:"SelfDeliveryFree"`
}
