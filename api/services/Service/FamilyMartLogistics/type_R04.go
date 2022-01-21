package FamilyMartLogistics

import (
	"encoding/xml"
)

/*
<?xml version="1.0" encoding="UTF-8"?>
<doc>
    <HEADER>
        <RDFMT>1</RDFMT>
        <SNCD>DFM</SNCD>
        <PRDT>2020-06-01</PRDT>
    </HEADER>
    <BODY>
        <R04>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <ShipmentNo>0001100340169</ShipmentNo>
            <DCReceiveDate>2020-06-01</DCReceiveDate>
            <DCReceiveStatus>1</DCReceiveStatus>
            <FlowType>N</FlowType>
        </R04>
    </BODY>
    <FOOTER>
        <RDFMT>3</RDFMT>
        <RDCNT>1</RDCNT>
    </FOOTER>
</doc>
*/

type FMLR04 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body  struct {
		Data []FMLR04Data `xml:"R04"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLR04Data struct {
	XMLName         xml.Name `xml:"R04"`
	RDFMT           string   `xml:"RDFMT"`
	ParentId        string   `xml:"ParentId"`
	EshopId         string   `xml:"EshopId"`
	ShipmentNo      string   `xml:"ShipmentNo"`
	DCReceiveDate   string   `xml:"DCReceiveDate"`
	DCReceiveStatus string   `xml:"DCReceiveStatus"`
	FlowType        string   `xml:"FlowType"`
}

func (receiver *FMLR04) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}
