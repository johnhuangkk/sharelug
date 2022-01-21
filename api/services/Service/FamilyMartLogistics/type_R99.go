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
        <PRDT>2020-06-01</PRDT>
    </HEADER>
    <BODY>
        <R99>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <ServiceType>1</ServiceType>
            <OPMode>1</OPMode>
            <ShipmentNo>0001100334462</ShipmentNo>
            <SPAmount>114</SPAmount>
            <StoreId>011926</StoreId>
            <SPAdate>20200601</SPAdate>
            <SPAstatus>1</SPAstatus>
            <SPFee>55</SPFee>
            <SPArate></SPArate>
        </R99>
        <R99>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <ServiceType>1</ServiceType>
            <OPMode>1</OPMode>
            <ShipmentNo>0002000032239</ShipmentNo>
            <SPAmount>49</SPAmount>
            <StoreId>011926</StoreId>
            <SPAdate>20200601</SPAdate>
            <SPAstatus>1</SPAstatus>
            <SPFee>0</SPFee>
            <SPArate></SPArate>
        </R99>
    </BODY>
    <FOOTER><RDFMT>3</RDFMT><RDCNT>2</RDCNT></FOOTER>
</doc>
*/

type FMLR99 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body   struct {
		Data []FMLR99Data `xml:"R99"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLR99Data struct {
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

func (receiver *FMLR99) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}
