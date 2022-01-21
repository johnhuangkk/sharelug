package MartHiLife

import "encoding/xml"

//<?xml version='1.0' encoding='UTF-8'?>
//<Doc>
//	<ShipmentNos>
//		<ParentId>124</ParentId>
//		<EshopId>901</EshopId>
//		<OrderNo>20RKB40416476</OrderNo>
//		<EcOrderNo>TSB30001968</EcOrderNo>
//		<OriginStoreId>3750</OriginStoreId>
//		<StoreType>2</StoreType>
//		<ChkMac>0C65AA5CC37A1D415AACD3E4C5D6DB32</ChkMac>
//	</ShipmentNos>
//</Doc>

type ChangeNotification struct {
	Doc struct{
		ParentId      string   `xml:"ParentId"`
		EshopId       string   `xml:"EshopId"`
		OrderNo       string   `xml:"OrderNo"`
		EcOrderNo     string   `xml:"EcOrderNo"`
		OriginStoreId string   `xml:"OriginStoreId"`
		StoreType     string   `xml:"StoreType"` //1:寄件店 2:取件店
		ChkMac        string   `xml:"ChkMac"`
	} `xml:"ShipmentNos"`
}

func (me *ChangeNotification) DecodeXML(data []byte) (err error) {
	return xml.Unmarshal(data,me)
}