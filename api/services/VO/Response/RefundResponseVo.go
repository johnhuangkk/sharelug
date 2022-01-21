package Response

type QueryRefundResponse struct {
	OrderId     string `json:"OrderId"`
	RefundedAmt int64  `json:"RefundedAmt"`
	OrderAmount int64  `json:"OrderAmount"`
	ProductAmt  int64  `json:"ProductAmt"`
	ShipFee     int64  `json:"ShipFee"`
}

type QueryRefundListResponse struct {
	RefundId         string `json:"RefundId"`
	RefundStatus     string `json:"RefundStatus"`
	RefundStatusText string `json:"RefundStatusText"`
	ApplyTime        string `json:"ApplyTime"`
	CompleteTime     string `json:"CompleteTime"`
	ThisRefund       int64  `json:"ThisRefund"`
	TotalRefund      int64  `json:"TotalRefund"`
}

type SearchRefundResponse struct {
	Tabs             RefundTabsResponse `json:"Tabs"`
	SearchRefundList []SearchRefundList `json:"SearchRefundList"`
}

type RefundTabsResponse struct {
	RefundAll  int64 `json:"RefundAll"`
	ReturnWait int64 `json:"ReturnWait"`
	ReturnSuc  int64 `json:"ReturnSuc"`
	RefundWait int64 `json:"RefundWait"`
	RefundSuc  int64 `json:"RefundSuc"`
}

type SearchRefundList struct {
	OrderId           string `json:"OrderId"`
	RefundType        string `json:"RefundType"`
	ReturnId          string `json:"ReturnId"`
	ReturnTime        string `json:"ReturnTime"`
	ReturnCheckTime   string `json:"ReturnCheckTime"`
	ReturnProductName string `json:"ReturnProductName"`
	ThisReturnQty     int64  `json:"ThisReturnQty"`
	TotalReturnQty    int64  `json:"TotalReturnQty"`
	ReturnStatus      string `json:"ReturnStatus"`
	ReturnStatusText  string `json:"ReturnStatusText"`
	RefundId          string `json:"RefundId"`
	RefundTime        string `json:"RefundTime"`
	ThisRefundAmt     int64  `json:"ThisRefundAmt"`
	TotalRefundAmt    int64  `json:"TotalRefundAmt"`
	RefundStatus      string `json:"RefundStatus"`
	RefundStatusText  string `json:"RefundStatusText"`
}
