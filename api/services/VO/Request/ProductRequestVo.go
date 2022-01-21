package Request

type NewProductParams struct {
	ProductImage    []string         `form:"ProductImage" json:"ProductImage"`                     //圖片
	ProductName     string           `form:"ProductName" json:"ProductName" validate:"required"`   //商品名稱
	ProductQty      int              `form:"ProductQty" json:"ProductQty"`                         //商品數量
	Price           string           `form:"Price" json:"Price" validate:"required"`               //價格
	ProductSpecList []NewProductSpec `form:"ProductSpecList" json:"ProductSpecList"`               //規格
	ShippingList    []NewShipping    `form:"ShippingList" json:"ShippingList" validate:"required"` //運送方式加運費
	PayWayList      []string         `form:"PayWayList" json:"PayWayList" validate:"required"`     //付款方式
	FormUrl         string           `form:"FormUrl" json:"FormUrl"`                               //Google表單連結
	ShipMerge       int              `form:"ShipMerge" json:"ShipMerge"`                           //是否合併運費
	IsSpec          int              `form:"IsSpec" json:"IsSpec"`                                 //是否有規格
	IsRealtime      int              `form:"IsRealtime" json:"IsRealtime"`                         //是否為即時帳單
	LimitKey        string           `form:"LimitKey" json:"LimitKey"`                             //限制
	LimitQty        int              `form:"LimitQty" json:"LimitQty"`                             //限制數量
	//IsFreeShip      bool             `form:"IsFreeShip" json:"IsFreeShip"`                         //是否免運費
}

type NewProductSpec struct {
	ProductSpecId string `form:"ProductSpecId"`
	ProductSpec   string `form:"ProductSpec"`
	Quantity      int    `form:"Quantity"`
}

type NewShipping struct {
	ShipType   string
	ShipFee    int
	ShipRemark string
}

type EditProductParams struct {
	ProductId       string           `form:"ProductId" json:"ProductId"`
	ProductImage    []string         `form:"ProductImage" json:"ProductImage"`                     //圖片
	ProductName     string           `form:"ProductName" json:"ProductName" validate:"required"`   //商品名稱
	ProductQty      int              `form:"ProductQty" json:"ProductQty"`                         //商品數量
	Price           string           `form:"Price" json:"Price" validate:"required"`               //價格
	ProductSpecList []NewProductSpec `form:"ProductSpecList" json:"ProductSpecList"`               //規格
	ShippingList    []NewShipping    `form:"ShippingList" json:"ShippingList" validate:"required"` //運送方式加運費
	PayWayList      []string         `form:"PayWayList" json:"PayWayList" validate:"required"`     //付款方式
	ShipMerge       int              `form:"ShipMerge" json:"ShipMerge"`                           //是否合併運費
	IsSpec          int              `form:"IsSpec" json:"IsSpec"`                                 //是否有規格
	IsRealtime      int              `form:"IsRealtime" json:"IsRealtime"`                         //是否為即時帳單
	StatusDown      int              `form:"StatusDown" json:"StatusDown"`
	StatusDelete    int              `form:"StatusDelete" json:"StatusDelete"`
	FormUrl         string           `form:"FormUrl" json:"FormUrl"` //Google表單連結
	LimitKey        string           `form:"LimitKey" json:"LimitKey"`                             //限制
	LimitQty        int              `form:"LimitQty" json:"LimitQty"`                             //限制數量
	//IsFreeShip      bool             `form:"IsFreeShip" json:"IsFreeShip"`                         //是否免運費
}
