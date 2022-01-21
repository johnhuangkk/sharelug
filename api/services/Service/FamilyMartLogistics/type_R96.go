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
        <R96>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <StoreId>003822</StoreId>
            <SPDate>2020-06-01</SPDate>
            <SPAmount>1280</SPAmount>
            <ServiceType>1</ServiceType>
            <ShipmentNo>800006430</ShipmentNo>
            <FlowType>R</FlowType>
        </R96>
        <R96>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <StoreId>008675</StoreId>
            <SPDate>2020-06-01</SPDate>
            <SPAmount>1980</SPAmount>
            <ServiceType>1</ServiceType>
            <ShipmentNo>800006451</ShipmentNo>
            <FlowType>N</FlowType>
        </R96>
        <R96>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <StoreId>009029</StoreId>
            <SPDate>2020-06-01</SPDate>
            <SPAmount>1280</SPAmount>
            <ServiceType>1</ServiceType>
            <ShipmentNo>800006445</ShipmentNo>
            <FlowType>R</FlowType>
        </R96>
        <R96>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <StoreId>010611</StoreId>
            <SPDate>2020-06-01</SPDate>
            <SPAmount>899</SPAmount>
            <ServiceType>1</ServiceType>
            <ShipmentNo>800006447</ShipmentNo>
            <FlowType>N</FlowType>
        </R96>
    </BODY>
    <FOOTER>
        <RDFMT>3</RDFMT>
        <RDCNT>4</RDCNT>
        <AMT>5439</AMT>
    </FOOTER>
</doc>
*/

type FMLR96 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body   struct {
		Data []FMLR96Data `xml:"R96"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLR96Data struct {
	RDFMT       string `xml:"RDFMT"`
	ParentId    string `xml:"ParentId"`
	EshopId     string `xml:"EshopId"`
	StoreId     string `xml:"StoreId"`
	SPDate      string `xml:"SPDate"`
	SPAmount    string `xml:"SPAmount"`
	ServiceType string `xml:"ServiceType"`
	ShipmentNo  string `xml:"ShipmentNo"`
	FlowType    string `xml:"FlowType"`
}

func (receiver *FMLR96) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}
