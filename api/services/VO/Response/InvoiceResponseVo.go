package Response

type InvoiceListResponse struct {
	InvoiceList  []InvoiceList `json:"InvoiceList"`
	InvoiceCount int64       `json:"InvoiceCount"`
}

type InvoiceList struct {
	OrderId           string `json:"OrderId"`           //訂單編號
	InvoiceStatus     string `json:"InvoiceStatus"`     //發票狀態
	InvoiceStatusText string `json:"InvoiceStatusText"` //發票狀態文字
	CarrierType       string `json:"CarrierType"`       //載具狀態
	CarrierTypeText   string `json:"CarrierTypeText"`   //載具狀態文字
	InvoiceNumber     string `json:"InvoiceNumber"`     //發票號碼
	Donate            string `json:"Donate"`            //捐贈單位
	CreateTime        string `json:"CreateTime"`        //開立時間
	CompanyBan        string `json:"CompanyBan"`        //統一編號
	CarrierId         string `json:"CarrierId"`         //載具編號
}
