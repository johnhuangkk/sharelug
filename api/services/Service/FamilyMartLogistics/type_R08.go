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
        <R08>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <ShipmentNo>M836000526207</ShipmentNo>
            <TotalAmount>1280</TotalAmount>
            <DCReturnDate>2020-06-01</DCReturnDate>
            <DCReturnStatus>1</DCReturnStatus>
            <FlowType>N</FlowType>
        </R08>
    </BODY>
    <FOOTER>
        <RDFMT>3</RDFMT>
        <RDCNT>1</RDCNT>
        <AMT>1280</AMT>
    </FOOTER>
</doc>
*/

type FMLR08 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body  struct {
		Data []FMLR08Data `xml:"R08"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLR08Data struct {
	RDFMT           string   `xml:"RDFMT"`
	ParentId        string   `xml:"ParentId"`
	EshopId         string   `xml:"EshopId"`
	ShipmentNo      string   `xml:"ShipmentNo"`
	TotalAmount     string   `xml:"TotalAmount"`
	DCReturnDate    string   `xml:"DCReturnDate"`
	DCReturnStatus  string   `xml:"DCReturnStatus"`
	FlowType        string   `xml:"FlowType"`
}

func (receiver *FMLR08) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}