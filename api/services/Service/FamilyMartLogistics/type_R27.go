package FamilyMartLogistics

import (
	"encoding/xml"
	"time"
)

/*
<?xml version="1.0" encoding="utf-8"?>
<doc>
    <HEADER>
        <RDFMT>1</RDFMT>
        <SNCD> DFM</SNCD>
        <PRDT>2020-06-01</PRDT>
    </HEADER>
    <BODY>
        <R27>
            <RDFMT>2</RDFMT>
            <ParentId>853</ParentId>
            <EshopId>0001</EshopId>
            <ShipmentNo>03918823876</ShipmentNo>
            <DCShipDate>2020-06-01</DCShipDate>
            <DCShipTime>141300</DCShipTime>
        </R27>
        <R27>
            <RDFMT>2</RDFMT>
            <ParentId>853</ParentId>
            <EshopId>0001</EshopId>
            <ShipmentNo>03918823876</ShipmentNo>
            <DCReceiveDate>2020-06-01</DCReceiveDate>
            <DCReceiveTime>141300</DCReceiveTime>
        </R27>
    </BODY>
    <FOOTER>
        <RDFMT>3</RDFMT>
        <RDCNT>2</RDCNT>
    </FOOTER>
</doc>
*/

type FMLR27 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body   struct {
		Data []FMLR27Data `xml:"R27"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLR27Data struct {
	RDFMT         string `xml:"RDFMT"`
	ParentId      string `xml:"ParentId"`
	EshopId       string `xml:"EshopId"`
	ShipmentNo    string `xml:"ShipmentNo"`
	DCShipDate    string `xml:"DCShipDate"`
	DCShipTime    string `xml:"DCShipTime"`
}

func (s *FMLR27Data) GetDateTime () string {
	t, _ := time.Parse(`2006-01-02150405`, s.DCShipDate + s.DCShipTime)
	return t.Format(`2006-01-02 15:04:05`)
}

func (receiver *FMLR27) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}
