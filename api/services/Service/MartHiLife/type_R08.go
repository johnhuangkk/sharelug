package MartHiLife

import "encoding/xml"

//<?xml version="1.0" encoding="utf-8"?>
//<doc>
//	<HEADER><RDFMT>1</RDFMT><SNCD>HIL</SNCD><PRDT>2017-10-06</PRDT></HEADER>
//	<BODY>
//		<R08>
//			<RDFMT>2</RDFMT>
//			<ParentId>124</ParentId>
//			<EshopId>901</EshopId>
//			<OrderNo>18Q1F41413172</OrderNo>
//			<TotalAmount>1660</TotalAmount>
//			<DCReturnDate>2017-10-05 09:00:05</DCReturnDate>
//			<DCReturnStatus>1</DCReturnStatus>
//			<FlowType>N</FlowType>
//		</R08>
//	</BODY>
//	<FOOTER><RDFMT>3</RDFMT><RDCNT>1</RDCNT></FOOTER>
//</doc>

type R08Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []R08Content `xml:"R08"`
	} `xml:"BODY"`
}

type R08Content struct {
	RDFMT            string `xml:"RDFMT"`
	ParentId         string `xml:"ParentId"`
	EshopId          string `xml:"EshopId"`
	OrderNo          string `xml:"OrderNo"`
	TotalAmount      string `xml:"TotalAmount"`
	DCReturnDate     string `xml:"DCReturnDate"`
	DCReceiveStatus  string `xml:"DCReceiveStatus"`
	FlowType         string `xml:"FlowType"`
}

func (receiver *R08Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
