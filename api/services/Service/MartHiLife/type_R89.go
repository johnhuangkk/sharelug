package MartHiLife

import "encoding/xml"

//<?xml version="1.0" encoding="utf-8"?>
//<doc>
//	<HEADER><RDFMT>1</RDFMT><SNCD>HIL</SNCD><PRDT>2011-10-06</PRDT></HEADER>
//	<BODY>
//		<R89>
//			<RDFMT>2</RDFMT>
//			<ParentId>124</ParentId>
//			<EshopId>901</EshopId>
//			<OrderNo>18Q1F41413172</OrderNo>
//			<OrderAmount>500</OrderAmount>
//			<ServiceType>1</ServiceType>
//			<OPMode>2</OPMode>
//			<StoreId>3750</StoreId>
//			<SPAdate>2017-10-05</SPAdate>
//			<SPAstatus>2</SPAstatus>
//			<SPFee>55</SPFee>
//			<SPArate></SPArate>
//		</R89>
//	</BODY>
//	<FOOTER><RDFMT>3</RDFMT><RDCNT>1</RDCNT></FOOTER>
//</doc>

type R89Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []R89Content `xml:"R89"`
	} `xml:"BODY"`
}

type R89Content struct {
	RDFMT       string `xml:"RDFMT"`
	ParentId    string `xml:"ParentId"`
	EshopId     string `xml:"EshopId"`
	OrderNo     string `xml:"OrderNo"`
	OrderAmount string `xml:"OrderAmount"`
	ServiceType string `xml:"ServiceType"`
	OPMode      string `xml:"OPMode"`
	StoreId     string `xml:"StoreId"`
	SPAdate     string `xml:"SPAdate"`
	SPAstatus   string `xml:"SPAstatus"`
	SPArate     string `xml:"SPArate"`
	SPFee       string `xml:"SPFee"`
}

func (receiver *R89Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
