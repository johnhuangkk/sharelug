package FamilyMartLogistics

import (
	"encoding/xml"
	"time"
)

/*
<?xml version="1.0" encoding="utf-8"?>
<doc>
    <HEADER>
        <RDFMT>1&gt;</RDFMT>
        <SNCD>DFM</SNCD>
        <PRDT>2020-06-01</PRDT>
    </HEADER>
    <BODY>
        <R29>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <OrderNo>03918823876</OrderNo>
            <EcOrderNo>12345678</EcOrderNo>
            <OrderDate>2020-06-15</OrderDate>
            <OrderTime>141300</OrderTime>
            <OPMode>1</OPMode>
            <StoreId>005989</StoreId>
            <StoreName>大溪介壽店</StoreName>
            <FlowType>N</FlowType>
            <StoreType>005989</StoreType>
        </R29>
        <R29>
            <RDFMT>2</RDFMT>
            <ParentId>861</ParentId>
            <EshopId>0001</EshopId>
            <OrderNo>03918823876</OrderNo>
            <EcOrderNo>12345678</EcOrderNo>
            <OrderDate>2020-06-15</OrderDate>
            <OrderTime>141300</OrderTime>
            <OPMode>1</OPMode>
            <StoreId>005989</StoreId>
            <StoreName>大溪介壽店</StoreName>
            <FlowType>N</FlowType>
            <StoreType>005989</StoreType>
        </R29>
    </BODY>
    <FOOTER>
        <RDFMT>3</RDFMT>
        <RDCNT>2</RDCNT>
    </FOOTER>
</doc>
*/

type FMLR29 struct {
	Header FMLDataHeader `xml:"HEADER"`
	Body  struct {
		Data []FMLR29Data `xml:"R29"`
	} `xml:"BODY"`
	Footer FMLDataFooter `xml:"FOOTER"`
}

type FMLR29Data struct {
	RDFMT     string `xml:"RDFMT"`
	ParentId  string `xml:"ParentId"`
	EshopId   string `xml:"EshopId"`
	OrderNo   string `xml:"OrderNo"`
	EcOrderNo string `xml:"EcOrderNo"`
	OrderDate string `xml:"OrderDate"`
	OrderTime string `xml:"OrderTime"`
	OPMode    string `xml:"OPMode"`
	StoreId   string `xml:"StoreId"`
	StoreName string `xml:"StoreName"`
	FlowType  string `xml:"FlowType"`
	StoreType string `xml:"StoreType"`
}

func (s *FMLR29Data) GetDateTime () string {
	t, _ := time.Parse(`2006-01-02150405`, s.OrderDate + s.OrderTime)
	return t.Format(`2006-01-02 15:04:05`)
}

func (receiver *FMLR29) DecodeXML(data []byte) error {
	return xml.Unmarshal(data, receiver)
}
