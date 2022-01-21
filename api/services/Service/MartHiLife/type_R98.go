package MartHiLife

import "encoding/xml"

//<?xml version="1.0" encoding="utf-8"?>
//<doc>
//	<HEADER><RDFMT>1</RDFMT><SNCD>HIL</SNCD><PRDT>2017-09-10</PRDT></HEADER>
//	<BODY>
//		<R98>
//			<RDFMT>2</RDFMT>
//			<ParentId>124</ParentId>
//			<EshopId>901</EshopId>
//			<OrderNo>18Q1F41413172</OrderNo>
//			<OPMode>2</OPMode>
//			<ServiceType>1</ServiceType>
//			<StoreId>3750</StoreId>
//			<SPAmount>0</SPAmount>
//			<SPAdate>2017-09-14</SPAdate>
//			<SPAstatus>1</SPAstatus>
//			<SPArate></SPArate>
//			<SPFee>60</SPFee>
//		</R98>
//	</BODY>
//	<FOOTER><RDFMT>3</RDFMT><RDCNT>1</RDCNT></FOOTER>
//</doc>

type R98Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []R98Content `xml:"R98"`
	} `xml:"BODY"`
}

type R98Content struct {
	RDFMT            string `xml:"RDFMT"`
	ParentId         string `xml:"ParentId"`
	EshopId          string `xml:"EshopId"`
	OrderNo          string `xml:"OrderNo"`
	OPMode   		 string `xml:"OPMode"`
	ServiceType      string `xml:"ServiceType"`
	StoreId  		 string `xml:"StoreId"`
	SPAmount         string `xml:"SPAmount"`
	SPAdate          string `xml:"SPAdate"`
	SPAstatus        string `xml:"SPAstatus"`
	SPArate          string `xml:"SPArate"`
	SPFee            string `xml:"SPFee"`
}

func (receiver *R98Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
