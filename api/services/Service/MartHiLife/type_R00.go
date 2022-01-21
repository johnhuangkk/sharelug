package MartHiLife

import "encoding/xml"

//<?xml version="1.0" encoding="utf-8"?>
//<doc>
//  <HEADER>
//		<RDFMT>1</RDFMT>
//		<SNCD>HIL</SNCD>
//		<PRDT>2019-10-03</PRDT>
//  </HEADER>
//  <BODY>
//    <R00>
//		<RDFMT>2</RDFMT>
//		<STOREID>4075</STOREID>
//		<STORE_NAME>龍潭大池店</STORE_NAME>
//		<MDC_START_DATE />
//		<MDC_END_DATE />
//		<ROUTER />
//		<STEP /> <STORE_ADDRESS>桃園市龍潭區神龍路234、236號</STORE_ADDRESS> <TEL_NO>032860264</TEL_NO>
//		<OLD_STORE />
//		<STORE_CLOSE_DATE />
//	  </R00>
//	  <R00>
//		<RDFMT>2</RDFMT><STOREID>4076</STOREID><STORE_NAME>嘉義國興店</STORE_NAME>
//		<MDC_START_DATE /><MDC_END_DATE /> <ROUTER /><STEP />
//		<STORE_ADDRESS>嘉義市東區興業東路160號</STORE_ADDRESS><TEL_NO>052949255</TEL_NO>
//		<OLD_STORE /><STORE_CLOSE_DATE />
//    </R00>
//  </BODY>
//  <FOOTER><RDFMT>3</RDFMT><RDCNT>2</RDCNT></FOOTER>
//</doc>

type R00Doc struct {
	XMLName xml.Name `xml:"doc"`
	Body    struct {
		Contents []R00Content `xml:"R00"`
	} `xml:"BODY"`
}

type F00DocHead struct {
	DocDate         string `xml:"DOCDATE"`
	FromPartnerCode string `xml:"FROMPARTNERCODE"`
	ToPartenerCode  string `xml:"TOPARTNERCODE"`
}

type R00Content struct {
	RDFMT          string `xml:"RDFMT"`
	StoreId        string `xml:"STOREID"`
	StoreName      string `xml:"STORE_NAME"`
	StoreAddress   string `xml:"STORE_ADDRESS"`
	TelNo          string `xml:"TEL_NO"`
	OldStore       string `xml:"OLD_STORE"`
	StoreCloseDate string `xml:"STORE_CLOSE_DATE"`
	MdcStareDate   string `xml:"MDC_START_DATE"`
	MdcEndDate     string `xml:"MDC_END_DATE"`
	ROUTE          string `xml:"ROUTER"`
	STEP           string `xml:"STEP"`
}

func (receiver *R00Doc) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
