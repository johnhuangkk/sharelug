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
        <RS4>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <ShipmentNo>800006429</ShipmentNo>
            <ReceiveStoreId>010307</ReceiveStoreId>
            <DCReceiveDate>2020-06-01</DCReceiveDate>
            <DCReceiveStatus>1</DCReceiveStatus>
            <FlowType>N</FlowType>
            <StoreType>2</StoreType>
            <StoreName><![CDATA[蘆竹大竹店]]></StoreName>
        </RS4>
    </BODY>
    <FOOTER><RDFMT>3</RDFMT><RDCNT>1</RDCNT></FOOTER>
</doc>
*/

type FMLRS4 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body   struct {
		Data []FMLRS4Data `xml:"RS4"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLRS4Data struct {
	RDFMT           string `xml:"RDFMT"`
	ParentId        string `xml:"ParentId"`
	EshopId         string `xml:"EshopId"`
	ShipmentNo      string `xml:"ShipmentNo"`
	ReceiveStoreId  string `xml:"ReceiveStoreId"`
	DCReceiveDate   string `xml:"DCReceiveDate"`
	DCReceiveStatus string `xml:"DCReceiveStatus"`
	FlowType        string `xml:"FlowType"`
	StoreType       string `xml:"StoreType"`
	StoreName       string `xml:"StoreName"`
}

func (receiver *FMLRS4) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}
