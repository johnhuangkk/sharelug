package FamilyMartLogistics

import (
	"encoding/xml"
)

/*
<?xml version="1.0" encoding="utf-8"?>
<doc>
    <HEADER>
        <RDFMT>1</RDFMT>
        <SNCD>DFM</SNCD>
        <PRDT>2020-06-15</PRDT>
    </HEADER>
    <BODY>
        <R98>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <ServiceType>1</ServiceType>
            <OPMode>1</OPMode>
            <ShipmentNo>0001100334462</ShipmentNo>
            <SPAmount>114</SPAmount>
            <StoreId>011926</StoreId>
            <SPAdate>2020-06-15</SPAdate>
            <SPAstatus>1</SPAstatus>
            <SPFee>60</SPFee>
            <SPArate></SPArate>
        </R98>
    </BODY>
    <FOOTER>
        <RDFMT>3</RDFMT>
        <RDCNT>1</RDCNT>
    </FOOTER>
</doc>
*/

type FMLR98 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body   struct {
		Data []FMLR98Data `xml:"R98"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLR98Data struct {
	RDFMT       string `xml:"RDFMT"`
	ParentId    string `xml:"ParentId"`
	EshopId     string `xml:"EshopId"`
	ServiceType string `xml:"ServiceType"`
	OPMode      string `xml:"OPMode"`
	ShipmentNo  string `xml:"ShipmentNo"`
	SPAmount    string `xml:"SPAmount"`
	StoreId     string `xml:"StoreId"`
	SPAdate     string `xml:"SPAdate"`
	SPAstatus   string `xml:"SPAstatus"`
	SPFee       string `xml:"SPFee"`
	SPArate     string `xml:"SPArate"`
}

func (receiver *FMLR98) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}
