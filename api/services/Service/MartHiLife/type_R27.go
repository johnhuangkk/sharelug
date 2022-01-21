package MartHiLife

import (
	"encoding/xml"
)

//<?xml version="1.0" encoding="utf-8"?>
//<doc>
//	<HEADER>
//		<RDFMT>1</RDFMT><SNCD>HIL</SNCD><PRDT>2019-01-19</PRDT>
//	</HEADER>
//	<BODY>
//		<R27>
//			<RDFMT>2</RDFMT>
//			<ParentId>124</ParentId>
//			<EshopId>901</EshopId>
//			<OrderNo>09DEH00504337</OrderNo>
//			<EcOrderNo>12411208137</EcOrderNo>
//			<SendStoreId>3750</SendStoreId>
//			<OrderDate>2019-01-18 00:49:36</OrderDate>
//			<SendStatus>1</SendStatus>
//		</R27>
//	</BODY>
//	<FOOTER>
//			<RDFMT>3</RDFMT>
//			<RDCNT>1</RDCNT>
//	</FOOTER>
//</doc>

type R27Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []R27Content `xml:"R27"`
	} `xml:"BODY"`
}

type R27Content struct {
	RDFMT       string `xml:"RDFMT"`
	ParentId    string `xml:"ParentId"`
	EshopId     string `xml:"EshopId"`
	OrderNo     string `xml:"OrderNo"`
	EcOrderNo   string `xml:"EcOrderNo"`
	SendStoreId string `xml:"SendStoreId"`
	OrderDate   string `xml:"OrderDate"`
	SendStatus  string `xml:"SendStatus"` //1:寄件成功 2:取消寄件
}

func (receiver *R27Doc) DecodeXML(data []byte) (err error) {
	return xml.Unmarshal(data, receiver)
}
