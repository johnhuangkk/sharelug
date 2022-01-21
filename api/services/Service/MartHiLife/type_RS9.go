package MartHiLife

import "encoding/xml"

//<?xml version="1.0" encoding="utf-8"?>
//<doc>
//	<HEADER>
//		<RDFMT>1</RDFMT><SNCD>HIL</SNCD><PRDT>2017-10-06</PRDT>
//	</HEADER>
//	<BODY>
//		<RS9>
//			<RDFMT>2</RDFMT>
//			<ParentId>124</ParentId>
//			<EshopId>901</EshopId>
//			<OrderNo>18Q1F41413172</OrderNo>
//			<ReceiveStoreId>3750</ReceiveStoreId>
//			<DCReceiveDate>2017-10-05</DCReceiveDate>
//			<DCReceiveStatus>2</DCReceiveStatus>
//			<StatusDetails>S03</StatusDetails>
//			<StatusRemark><![CDATA[小物流遺失]]></StatusRemark>
//			<FlowType>N</FlowType>
//		</RS9>
//	</BODY>
//	<FOOTER>
//		<RDFMT>3</RDFMT><RDCNT>1</RDCNT>
//	</FOOTER>
//</doc>

type RS9Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []RS9Content `xml:"RS9"`
	} `xml:"BODY"`
}

type RS9Content struct {
	RDFMT           string `xml:"RDFMT"`
	ParentId        string `xml:"ParentId"`
	EshopId         string `xml:"EshopId"`
	OrderNo         string `xml:"OrderNo"`
	ReceiveStoreId  string `xml:"ReceiveStoreId"`
	DCReceiveDate   string `xml:"DCReceiveDate"`
	DCReceiveStatus string `xml:"DCReceiveStatus"`
	StatusDetails   string `xml:"StatusDetails"`
	StatusRemark    string `xml:"StatusRemark"`
	FlowType        string `xml:"FlowType"`
}

func (receiver *RS9Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
