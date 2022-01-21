package FamilyMartLogistics

import (
	"encoding/xml"
)

type Map map[string]interface{}

type FMLData struct {
	XMLName xml.Name      `xml:"doc"`
	Header  FMLDataHeader `xml:"HEADER"`
	Body    FMLDataBody   `xml:"BODY"`
	Footer  FMLDataFooter `xml:"FOOTER"`
}

type FMLDataHeader struct {
	RDFMT string `xml:"RDFMT"` // 區別碼
	SNCD  string `xml:"SNCD"`
	PRDT  string `xml:"PRDT"`
}

type FMLDataBody struct {
	R22Body Map `xml:"R22"`
}

type FMLDataFooter struct {
	RDFMT string `xml:"RDFMT"`
	RDCNT string `xml:"RDCNT"`
	AMT   string `xml:"AMT"`
}