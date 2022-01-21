package FamilyMartLogistics

import "encoding/xml"

//<?xml version="1.0" encoding="UTF-8"?>
//<Doc>
//	<ShipmentNos>
//		<ParentId>888</ParentId>
//		<EshopId>0001</EshopId>
//		<OrderNo>04901154451</OrderNo>
//		<EcOrderNo>1234567890</EcOrderNo>
//		<StoreType>2</StoreType>
//		<ErrorCode>000</ErrorCode>
//		<ErrorMessage>成功</ErrorMessage>
//	</ShipmentNos>
//	<ShipmentNos>
//		<ParentId>888</ParentId>
//		<EshopId>0001</EshopId>
//		<OrderNo>04901154661</OrderNo>
//		<EcOrderNo>1234567660</EcOrderNo>
//		<StoreType>2</StoreType>
//		<ErrorCode>000</ErrorCode>
//		<ErrorMessage>成功</ErrorMessage>
//	</ShipmentNos>
//</Doc>

type OrderSwitchResponse struct {
	Content []OrderSwitchResponseObject `xml:"ShipmentNos"`
}

type OrderSwitchResponseObject struct {
	ParentId     string `xml:"ParentId"`
	EshopId      string `xml:"EshopId"`
	OrderNo      string `xml:"OrderNo"`
	EcOrderNo    string `xml:"EcOrderNo"`
	StoreType    string `xml:"StoreType"`
	ErrorCode    string `xml:"ErrorCode"`
	ErrorMessage string `xml:"ErrorMessage"`
}

func (receiver *OrderSwitchResponse) DecodeXML(data []byte) error {
	return xml.Unmarshal(data,&receiver)
}