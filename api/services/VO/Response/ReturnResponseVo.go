package Response

type QueryReturnResponse struct {
	OrderId     string              `json:"OrderId"`
	ProductList []ReturnProductList `json:"ProductList"`
}

type ReturnProductList struct {
	ProductName     string `json:"ProductName"`
	ProductSpecId   string `json:"ProductSpecId"`
	ProductSpecName string `json:"ProductSpecName"`
	ProductPrice    int64  `json:"ProductPrice"`
	ProductQty      int64  `json:"ProductQty"`
	Refundable      int64  `json:"Refundable"`
}

type QueryReturnListResponse struct {
	ReturnId         string `json:"ReturnId"`
	ReturnStatus     string `json:"ReturnStatus"`
	ReturnStatusText string `json:"ReturnStatusText"`
	ReturnTime       string `json:"ReturnTime"`
	CompleteTime     string `json:"CompleteTime"`
	ProductName      string `json:"ProductName"`
	ThisReturn       int64  `json:"ThisReturn"`
	TotalReturn      int64  `json:"TotalReturn"`
}

type RealTimesExtensionResponse struct {
	Time  string `json:"Time"`
}