package FamilyMartLogistics

import (
	"encoding/xml"
	"log"
)

//<?xml version="1.0" encoding="UTF-8"?>
//<Doc>
//	<ParentId>888</ParentId>
//	<EshopId>0001</EshopId>
//	<OrderService>
//		<OrderNo>06900000001</OrderNo>
//	</OrderService>
//	<OrderService>
//		<OrderNo>06900000011</OrderNo>
//	</OrderService>
//</Doc>
//

type OrdersPrintRequest struct {
	XMLName	xml.Name `xml:"Doc"`
	ParentId        string
	EshopId         string
	Orders          []OrdersPrintRequestOrder
}

type OrdersPrintRequestOrder struct {
	//XMLName	xml.Name `xml:"Doc"`
	OrderNo          string
}

func (receiver *OrdersPrintRequest) EncodeXML() (string,bool) {
	data,err := xml.Marshal(receiver)
	if err != nil {
		log.Println("EncodeXML:",err)
		return "",false
	}
	return xml.Header+string(data),true
}