package HiLifeMart

import "encoding/xml"

type HiLifParams struct {
	Data      string `form:"Data" json:"Data"`
}

type Response struct {
	XMLName     xml.Name `xml:"Doc"`
	ShipmentNos []ResponseBody
}

type ResponseBody struct {
	ParentId     string `xml:"ParentId" json:"ParentId"`
	EshopId      string `xml:"EshopId" json:"EshopId"`
	OrderNo      string `xml:"OrderNo" json:"OrderNo"`
	EcOrderNo    string `xml:"EcOrderNo" json:"EcOrderNo"`
	ErrorCode    string `xml:"ErrorCode" json:"ErrorCode"`
	ErrorMessage string `xml:"ErrorMessage" json:"ErrorMessage"`
}


//// 閉轉通知回覆給萊爾富
func (r *ResponseBody) RspSwitch(s SwitchBody, code, message string)  {
	r.ParentId = s.ParentId
	r.EshopId = s.EshopId
	r.OrderNo = s.OrderNo
	r.EcOrderNo = s.EcOrderNo
	r.ErrorCode = code
	r.ErrorMessage = message
}

type SwitchBody struct {
	ParentId      string `xml:"ParentId" json:"ParentId"`           // 母代碼  length 3
	EshopId       string `xml:"EshopId" json:"EshopId"`             // 子廠商代號 length 3
	OrderNo       string `xml:"OrderNo" json:"OrderNo"`             // 寄件代碼 length 13
	EcOrderNo     string `xml:"EcOrderNo" json:"EcOrderNo"`         // 訂單單號 length 11 萊爾富自己產生的單號
	OriginStoreId string `xml:"OriginStoreId" json:"OriginStoreId"` // 閉店原始店號 length 4
	StoreType     string `xml:"StoreType" json:"StoreType"`         // 閉店原始店號 1:寄件店 2:取件店  length 1
	ChkMac        string `xml:"ChkMac" json:"ChkMac"`               // 檢查碼
}

// 閉轉店即時通知
type Switch struct {
	ShipmentNos []SwitchBody
}
