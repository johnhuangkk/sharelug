package Response

type CartPayResponse struct {
	Product      CartsProductModel `json:"ProductList"`
	Subtotal     float64           `json:"Subtotal"`
	Shipping     string            `json:"Shipping"`
	ShipFee      float64           `json:"ShipFee"`
	Coupon       float64           `json:"Coupon"`
	CouponNumber string            `json:"CouponNumber"`
	BeforeTotal  float64           `json:"BeforeTotal"`
	Total        float64           `json:"Total"`
	IsCoupon     bool              `json:"IsCoupon"`
	ShipList     []ShippingMode    `json:"ShipList"`
	PayWayList   []PayWayMode      `json:"PayWayList"`
	StoreId      string            `json:"StoreId"`
}

type CartsProductModel struct {
	Merge      []CartsProductList `json:"CanNotMergeList"`
	MergeFee   int                `json:"CanNotMergeFee"`
	General    []CartsProductList `json:"GeneralList"`
	GeneralFee int                `json:"GeneralFee"`
	Free       []CartsProductList `json:"ShipFreeList"`
	FreeFee    int                `json:"ShipFreeFee"`
}

type CartsProductList struct {
	ProductId        string       `json:"ProductId"`
	ProductSpecId    string       `json:"ProductSpecId"`
	ProductName      string       `json:"ProductName"`
	ProductSpec      string       `json:"ProductSpec"`
	ProductImage     string       `json:"ProductImage"`
	Quantity         int          `json:"Quantity"`
	Price            int          `json:"Price"`
	ShipMerge        int          `json:"ShipMerge"`
	ShipMode         []ShipMode   `json:"ShipList"`
	PayWayMode       []PayWayMode `json:"PayWayList"`
	ShipFee          int          `json:"ShipFee"`
	FormUrl          string       `json:"FormUrl"`
	LimitKey         string       `json:"LimitKey"`
	LimitQty         int          `json:"LimitQty"`
	FreeShipKey      string       `json:"FreeShipKey"`
	FreeShip         int64        `json:"FreeShip"`
	SelfDeliveryKey  string       `json:"SelfDeliveryKey"`
	SelfDeliveryFree int64        `json:"SelfDeliveryFree"`
	ShipList         string       `json:"-"`
	PayWayList       string       `json:"-"`
}

type ShippingMode struct {
	Type   string `json:"Type"`
	Text   string `json:"Text"`
	Price  int    `json:"Price"`
	Remark string `json:"Remark"`
}

type ShipMode struct {
	Type string `json:"Type"`
	Text string `json:"Text"`
}

type PayWayMode struct {
	Type string `json:"Type"`
	Text string `json:"Text"`
}

type StoreProduct struct {
	ProductId         string             `json:"ProductId"`
	ProductName       string             `json:"ProductName"`
	ProductImage      []string           `json:"ProductImageList"`
	ProductPrice      int64              `json:"ProductPrice"`
	ProductIsSpec     int                `json:"ProductIsSpec"`
	ProductStatus     string             `json:"ProductStatus"`
	ProductStatusText string             `json:"ProductStatusText"`
	ProductExpireTime string             `json:"ProductExpireTime"`
	ProductCancelTime string             `json:"ProductCancelTime"`
	ProductShipMerge  int                `json:"ProductShipMerge"`
	ProductQrcode     string             `json:"ProductQrcode"`
	TotalStock        int64              `json:"TotalStock"`
	FormUrl           string             `json:"FormUrl"`
	LimitKey          string             `json:"LimitKey"`
	LimitQty          int                `json:"LimitQty"`
	FreeShipKey       string             `json:"FreeShipKey"`
	FreeShip          int64              `json:"FreeShip"`
	SelfDeliveryKey   string             `json:"SelfDeliveryKey"`
	SelfDeliveryFree  int64              `json:"SelfDeliveryFree"`
	ProductSpecList   []StoreProductSpec `json:"ProductSpecList"`
	ProductShipList   []ShippingMode     `json:"ProductShipList"`
	ProductPayWayList []PayWayMode       `json:"ProductPayWayList"`
}

type StoreProductSpec struct {
	ProductSpecId string `json:"ProductSpecId"`
	Spec          string `json:"Spec"`
	Quantity      int64  `json:"Quantity"`
	Price         int64  `json:"Price"`
}

type StoreProductsResponse struct {
	StoreId           string                   `json:"StoreId"`
	StoreName         string                   `json:"StoreName"`
	StoreImage        string                   `json:"StoreImage"`
	StoreStatus       string                   `json:"StoreStatus"`
	StoreSocialMedia  StoreSocialMediaResponse `json:"StoreSocialMedia"`
	FreeShipKey       string                   `json:"FreeShipKey"`
	FreeShip          int64                    `json:"FreeShip"`
	SelfDeliveryKey   string                   `json:"SelfDeliveryKey"`
	SelfDeliveryFree  int64                    `json:"SelfDeliveryFree"`
	VerifyIdentity    int64                    `json:"VerifyIdentity"`
	ProductCount      int64                    `json:"ProductCount"`
	ProductTotal      int64                    `json:"ProductTotal"`
	ProductSellCount  int64                    `json:"ProductSellCount"`
	ProductDownCount  int64                    `json:"ProductDownCount"`
	ProductStockCount int64                    `json:"ProductStockCount"`
	ProductList       []StoreProduct           `json:"ProductList"`
}

type ProductResponse struct {
	ProductId        string                   `json:"ProductId"`
	ProductName      string                   `json:"ProductName"`
	ProductImageList []string                 `json:"ProductImageList"`
	StorePicture     string                   `json:"StorePicture"`
	StoreName        string                   `json:"StoreName"`
	StoreSocialMedia StoreSocialMediaResponse `json:"StoreSocialMedia"`
	StoreId          string                   `json:"StoreId"`
	ShortUrl         string                   `json:"ShortUrl"`
	FullUrl          string                   `json:"FullUrl"`
	QRCode           string                   `json:"QrCode"`
	Price            int64                    `json:"Price"`
	TotalStock       int64                    `json:"TotalStock"`
	IsSpec           int                      `json:"IsSpec"`
	IsRealTime       int                      `json:"IsRealTime"`
	ExpireDate       string                   `json:"ExpireDate"`
	ShipMerge        int                      `json:"ShipMerge"`
	FormUrl          string                   `json:"FormUrl"`
	LimitKey         string                   `json:"LimitKey"`
	LimitQty         int                      `json:"LimitQty"`
	FreeShipKey      string                   `json:"FreeShipKey"`
	FreeShip         int64                    `json:"FreeShip"`
	SelfDeliveryKey  string                   `json:"SelfDeliveryKey"`
	SelfDeliveryFree int64                    `json:"SelfDeliveryFree"`
	Status           string                   `json:"Status"`
	UpdateTime       string                   `json:"UpdateTime"`
	ProductSpecList  []ProductSpec            `json:"ProductSpecList"`
	ShippingList     []ShippingMode           `json:"ShippingList"`
	ProductPayWay    []PayWayMode             `json:"ProductPayWayList"`
}

type ProductSpec struct {
	ProductSpecId string `json:"ProductSpecId"`
	Spec          string `json:"Spec"`
	Quantity      int64  `json:"Quantity"`
	Price         int64  `json:"Price"`
}

type CartsCountResponse struct {
	Count int `json:"Count"`
}

type StoreRealTimesResponse struct {
	Tabs         StoreRealTimeTabs `json:"Tabs"`
	ProductCount int64             `json:"ProductCount"`
	ProductList  []StoreProduct    `json:"ProductList"`
}

type StoreRealTimeTabs struct {
	ValidBill   int64 `json:"ValidBill"`
	CancelBill  int64 `json:"CancelBill"`
	OverdueBill int64 `json:"OverdueBill"`
}

type MemberInfoResponse struct {
	Balance    int64                `json:"Balance"`
	MemberCard []MemberCardResponse `json:"MemberCard"`
}

type MemberCardResponse struct {
	CardId      string `json:"CardId"`
	CardNumber  string `json:"CardNumber"`
	ExpiryDate  string `json:"ExpiryDate"`
	DefaultCard bool   `json:"DefaultCard"`
}

type AddressResponse struct {
	Username string `json:"Username"`
	Address  string `json:"Address"`
}
