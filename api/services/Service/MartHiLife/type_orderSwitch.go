package MartHiLife

import (
	"encoding/xml"
	"github.com/spf13/viper"
)

type OrderSwitchRequest struct {
	ParentId        string
	EshopId         string
	EcDcNo          string
	EcCvs           string
	ShipNo 			string
	EcOrderNo       string
	RcvStoreId      string
	StoreType       string
}

func (os *OrderSwitchRequest) GetRequest ()  {
	config := viper.GetStringMapString(`MartHiLife`)

	os.ParentId = config[`parentid`]
	os.EshopId = config[`eshopid`]
	os.EcDcNo = config[`ecdcno`]
	os.EcCvs = config[`eccvs`]
}

//<?xml version="1.0" encoding="utf-8"?>
//<Doc>
//	<ShipmentNos>
//		<ParentId>124</ParentId>
//		<EshopId>901</EshopId>
//		<OrderNo>18GFN12345678</OrderNo>
//		<EcOrderNo>12400000001</EcOrderNo>
//		<ErrorCode>000</ErrorCode>
//		<ErrorMessage>成功</ErrorMessage>
//  </ShipmentNos>
//</Doc>

type OrderSwitchResponse struct {
	Doc struct{
		ParentId        string `xml:"ParentId"`
		EshopId         string `xml:"EshopId"`
		ShipNo 			string `xml:"OrderNo"`
		EcOrderNo       string `xml:"EcOrderNo"`
		ErrorCode       string `xml:"ErrorCode"`
		ErrorMessage    string `xml:"ErrorMessage"`
	} `xml:"ShipmentNos"`
}

func (me *OrderSwitchResponse)	DecodeXML(data []byte) (err error) {
	return xml.Unmarshal(data, me)
}