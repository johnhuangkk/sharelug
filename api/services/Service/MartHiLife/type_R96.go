package MartHiLife

import "encoding/xml"

//<?xml version="1.0" encoding="utf-8"?>
//<doc>
//	<HEADER><RDFMT>1</RDFMT><SNCD>HIL</SNCD><PRDT>2017-10-06</PRDT></HEADER>
//	<BODY>
//		<R96>
//			<RDFMT>2</RDFMT>
//			<ParentId>124</ParentId>
//			<EshopId>901</EshopId>
//			<StoreId>3750</StoreId>
//			<SPDate>2017-10-05 09:00:05</SPDate>
//			<SPAmount>65</SPAmount>
//			<ServiceType>1</ServiceType>
//			<OrderNo>18Q1F41413172</OrderNo>
//			<FlowType>N</FlowType>
//		</R96>
//	</BODY>
//	<FOOTER><RDFMT>3</RDFMT><RDCNT>1</RDCNT><AMT>65</AMT></FOOTER>
//</doc>

type R96Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []R96Content `xml:"R96"`
	} `xml:"BODY"`
}

type R96Content struct {
	RDFMT            string `xml:"RDFMT"`
	ParentId         string `xml:"ParentId"`
	EshopId          string `xml:"EshopId"`
	StoreId          string `xml:"StoreId"`
	SPDate           string `xml:"SPDate"`
	SPAmount         string `xml:"SPAmount"`
	ServiceType      string `xml:"ServiceType"`
	OrderNo          string `xml:"OrderNo"`
	FlowType         string `xml:"FlowType"`
}

func (receiver *R96Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
