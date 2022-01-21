package MartHiLife

import "encoding/xml"

//<?xml version="1.0" encoding="utf-8"?> <doc>
//<HEADER>
//	<RDFMT>1</RDFMT>
//	<SNCD>HIL</SNCD> <PRDT>2017-10-06</PRDT>
//</HEADER>
//<BODY>
//	<R04>
//		<RDFMT>2</RDFMT>
//		<ParentId>124</ParentId>
//		<EshopId>901</EshopId>
//		<OrderNo>18Q1F41413172</OrderNo>
//		<DCReceiveDate>2017-10-06 09:00:05</DCReceiveDate>
//		<DCReceiveStatus>1</DCReceiveStatus>
//		<FlowType>N</FlowType>
//	</R04>
//</BODY>
//<FOOTER><RDFMT>3</RDFMT><RDCNT>1</RDCNT></FOOTER>
//</doc>

type R04Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []R04Content `xml:"R04"`
	} `xml:"BODY"`
}

type R04Content struct {
	RDFMT           string `xml:"RDFMT"`
	ParentId        string `xml:"ParentId"`
	EshopId         string `xml:"EshopId"`
	OrderNo         string `xml:"OrderNo"`
	DCReceiveDate   string `xml:"DCReceiveDate"`
	DCReceiveStatus string `xml:"DCReceiveStatus"`
	FlowType        string `xml:"FlowType"`
}

func (receiver *R04Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
