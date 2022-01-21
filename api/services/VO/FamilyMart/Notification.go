package FamilyMart

import (
	"encoding/xml"
	"net/url"
	"time"
)

type FamilyParams struct {
	ApiKey    string `form:"ApiKey" json:"ApiKey"`
	Data      string `form:"Data" json:"Data"`
	TimeStamp string `form:"TimeStamp" json:"TimeStamp"`
	Signature string `form:"Signature" json:"Signature"`
}

func (e *FamilyParams) DataUrlEncode() string {
	data, _ := url.QueryUnescape(e.Data)
	return data
}

/**
	全家即時通知共用
 */
type ShipmentNos struct {
	ParentId      string `xml:"ParentId" json:"ParentId"` // 母代碼 length 3
	EshopId       string `xml:"EshopId" json:"EshopId"` // 子廠商代號 length 4
	OrderNo       string `xml:"OrderNo" json:"OrderNo"` // 訂單編號 length 11
	EcOrderNo     string `xml:"EcOrderNo" json:"EcOrderNo"` // 廠商訂單編號 length 13
	OrderDate     string `xml:"OrderDate" json:"OrderDate"` // 交易日期 length 10 YYYY-MM-DD
	OrderTime     string `xml:"OrderTime" json:"OrderTime"` // 交易時間 length 6 HH24MMSS
}

// 取得時間
func (s *ShipmentNos) GetDateTime () string {
	t, _ := time.Parse(`2006-01-02150405`, s.OrderDate + s.OrderTime)
	return t.Format(`2006-01-02 15:04:05`)
}

/**
寄件通知 SEND_SENTNOTIFY_API
寄件離店即時通知 SEND_LEAVEANOTIFY_API
*/
type Send struct {
	ShipmentNos
	SendStoreType string `xml:"SendStoreType" json:"SendStoreType"` // 寄件通路 length 1
	SendStoreId   string `xml:"SendStoreId" json:"SendStoreId"` // 寄件店號 length 6
}

// 進店即時通知
type Enter struct {
	Pickup
	StoreType string `xml:"StoreType" json:"StoreType"` // 1:寄件店 2:取件店 3:轉換店
}

// 取件通知
type Pickup struct {
	ShipmentNos
	RcvStoreType string `xml:"RcvStoreType" json:"RcvStoreType"`
	RcvStoreId string `xml:"RcvStoreId" json:"RcvStoreId"` //現行店號
	FlowType string `xml:"FlowType" json:"FlowType"` // N:進貨 R:退貨
}

// 閉轉店即時通知
type Switch struct {
	ParentId      string `xml:"ParentId" json:"ParentId"` // 母代碼 length 3
	EshopId       string `xml:"EshopId" json:"EshopId"` // 子廠商代號 length 4
	OrderNo       string `xml:"OrderNo" json:"OrderNo"` // 訂單編號 length 11
	EcOrderNo     string `xml:"EcOrderNo" json:"EcOrderNo"` // 廠商訂單編號 length 13
	//ShipmentNos
	OriginStoreType     string `xml:"OriginStoreType" json:"OriginStoreType"` // 閉店原始通路 length 1
	OriginStoreId     string `xml:"OriginStoreId" json:"OriginStoreId"` // 閉店原始店號 length 6
	StoreType     string `xml:"StoreType" json:"StoreType"` // 閉店原始店號 1:寄件店 2:取件店  length 1
}

// SEND_SENT_NOTIFY_API  SEND_LEAVEANOTIFY_API
type SendDoc struct {
	ShipmentNos []Send
}

type EnterDoc struct {
	ShipmentNos []Enter
}

type PickupDoc struct {
	ShipmentNos []Pickup
}

type SwitchDoc struct {
	ShipmentNos []Switch
}

type ResponseXml struct {
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
	StoreType    string `xml:"StoreType,omitempty" json:"StoreType"`
}

type ErrorInfo struct {
	ErrorCode    string
	ErrorMessage string
}

// 設定錯誤訊息
func (rb *ResponseBody) SetErrorInfo(errInfo ErrorInfo, s ShipmentNos)  {
	rb.ErrorMessage = errInfo.ErrorMessage
	rb.ErrorCode = errInfo.ErrorCode
	rb.EcOrderNo = s.EcOrderNo
	rb.ParentId = s.ParentId
	rb.EshopId = s.EshopId
	rb.OrderNo = s.OrderNo
}


func (rb *ResponseBody) SetSwitchErrorInfo(errInfo ErrorInfo, s Switch)  {
	rb.ErrorMessage = errInfo.ErrorMessage
	rb.ErrorCode = errInfo.ErrorCode
	rb.EcOrderNo = s.EcOrderNo
	rb.ParentId = s.ParentId
	rb.EshopId = s.EshopId
	rb.OrderNo = s.OrderNo
}