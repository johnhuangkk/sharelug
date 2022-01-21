package FamilyMartLogistics

import (
	"encoding/xml"
)

//<?xml version="1.0" encoding="UTF-8"?> <doc>
//<HEADER><RDFMT>1</RDFMT> <SNCD>861</SNCD> <PRDT>2020-06-15</PRDT></HEADER>
//<BODY>
//	<I00>
//		<RDFMT>2</RDFMT>
//		<StoreId>011926</StoreId>
//		<StoreName>新營東興店</StoreName>
//		<MdcStareDate>2018-12-15</MdcStareDate> <MdcEndDate/>
//		<ROUTE>BA</ROUTE>
//		<STEP>017</STEP>
//		<StoreAddress>台南市新營區新東里東興路188號1樓</StoreAddress>
//		<TelNo/>
//		<OldStoreE>008062</OldStore>
//		<StoreCloseDate/>
//      <Area></Area>
//		<EquipmentID></EquipmentID>
//  </I00>
//</BODY>

type FMLI00 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body   struct {
		Data []FMLI00Data `xml:"I00"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLI00Data struct {
	RDFMT          string `xml:"RDFMT"`
	StoreId        string `xml:"StoreId"`
	StoreName      string `xml:"StoreName"`
	MdcStareDate   string `xml:"MdcStareDate"`
	MdcEndDate     string `xml:"MdcEndDate"`
	ROUTE          string `xml:"ROUTE"`
	STEP           string `xml:"STEP"`
	StoreAddress   string `xml:"StoreAddress"`
	TelNo          string `xml:"TelNo"`
	OldStore       string `xml:"OldStore"`
	StoreCloseDate string `xml:"StoreCloseDate"`
	Area           string `xml:"Area"`
	EquipmentID    string `xml:"EquipmentID"`
}


func (receiver *FMLI00) DecodeXML(data []byte) error {
	err := xml.Unmarshal(data, receiver)
	if err != nil {
		return err
	}
	return nil
}