package MartHiLife

import "encoding/xml"

//<?xml version=”1.0” encoding=”utf-8”?>
//<doc>
//	<HEADER><RDFMT>1</RDFMT><SNCD>HIL</SNCD><PRDT>2017-10-06</PRDT></HEADER>
//	<BODY>
//		<R29>
//			<RDFMT>2</RDFMT>
//			<ParentId>124</ParentId>
//			<EshopId>901</EshopId>
//			<OrderNo>18Q1F41413172</OrderNo>
//			<StoreId>3750</StoreId>
//			<SPDate>2017-10-05 09:00:05</SPDate>
//			<SPAmount>65</SPAmount>
//			<FlowType>N</FlowType>
//			<ServiceType>1</ServiceType>
//		</R29>
//	</BODY>
//	<FOOTER><RDFMT>3</RDFMT><RDCNT>1</RDCNT><AMT>65</AMT></FOOTER>
//</doc>

type R29Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []R29Content `xml:"R29"`
	} `xml:"BODY"`
}

type R29Content struct {
	RDFMT            string `xml:"RDFMT"`
	ParentId         string `xml:"ParentId"`
	EshopId          string `xml:"EshopId"`
	OrderNo          string `xml:"OrderNo"`
	StoreId          string `xml:"StoreId"`
	SPDate           string `xml:"SPDate"`
	SPAmount         string `xml:"SPAmount"`
	FlowType         string `xml:"FlowType"`
	ServiceType      string `xml:"ServiceType"`
}

func (receiver *R29Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
