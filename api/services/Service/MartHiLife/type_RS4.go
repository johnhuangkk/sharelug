package MartHiLife

import "encoding/xml"

//<?xml version="1.0" encoding="utf-8"?>
//<doc>
//	<HEADER>
//		<RDFMT>1</RDFMT><SNCD>HIL</SNCD><PRDT>2017-10-06</PRDT>
//	</HEADER>
//	<BODY>
//		<RS4>
//			<RDFMT>2</RDFMT>
//			<ParentId>124</ParentId>
//			<EshopId>901</EshopId>
//			<OrderNo>18Q1F41413172</OrderNo>
//			<ReceiveStoreId>3750</ReceiveStoreId>
//			<DCReceiveDate>2017-10-05 09:00:05</DCReceiveDate>
//			<DCReceiveStatus>1</DCReceiveStatus>
//			<FlowType>N</FlowType>
//			<StoreType>2</StoreType>
//			<StoreName><![CDATA[板橋站前店]]></StoreName>
//		</RS4>
//	</BODY>
//	<FOOTER>
//		<RDFMT>3</RDFMT><RDCNT>1</RDCNT>
//	</FOOTER>
//</doc>

type RS4Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []RS4Content `xml:"RS4"`
	} `xml:"BODY"`
}

type RS4Content struct {
	RDFMT            string `xml:"RDFMT"`
	ParentId         string `xml:"ParentId"`
	EshopId          string `xml:"EshopId"`
	OrderNo          string `xml:"OrderNo"`
	ReceiveStoreId   string `xml:"ReceiveStoreId"`
	DCReceiveDate    string `xml:"DCReceiveDate"`
	DCReceiveStatus  string `xml:"DCReceiveStatus"`
	FlowType         string `xml:"FlowType"`
	StoreType        string `xml:"StoreType"`
	StoreName        string `xml:"StoreName"`
}

func (receiver *RS4Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
