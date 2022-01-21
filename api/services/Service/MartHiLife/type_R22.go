package MartHiLife

import "encoding/xml"

//<?xml version="1.0" encoding="utf-8"?>
//<doc>
//<HEADER>
//	<RDFMT>1</RDFMT>
//	<SNCD>HIL</SNCD>
//	<PRDT>2017-10-06</PRDT>
//</HEADER>
//<BODY>
//	<R22>
//		<RDFMT>2</RDFMT>
//		<ParentId>124</ParentId>
//		<EshopId>901</EshopId>
//		<OrderNo>18Q1F414131721</OrderNo>
//		<EcOrderNo>12400000001</EcOrderNo>
//		<SendStoreId>3750</SendStoreId>
//		<OrderDate>2017-10-05 00:49:36</OrderDate>
//		<SendStatus>1</SendStatus>
//	</R22>
//</BODY>
//<FOOTER>
//	<RDFMT>3</RDFMT>
//	<RDCNT>1</RDCNT>
//</FOOTER>
//</doc>

type R22Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []R22Content `xml:"R22"`
	} `xml:"BODY"`
}

type R22Content struct {
	RDFMT       string `xml:"RDFMT"`
	ParentId    string `xml:"ParentId"`
	EshopId     string `xml:"EshopId"`
	OrderNo     string `xml:"OrderNo"`
	EcOrderNo   string `xml:"EcOrderNo"`
	SendStoreId string `xml:"SendStoreId"`
	OrderDate   string `xml:"OrderDate"`
	SendStatus  string `xml:"SendStatus"`
}

func (receiver *R22Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
