package OKMart

import (
	"encoding/xml"
)

//<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
//  <soap:Body>
//    <ORDERS_ADDResponse xmlns="http://tempuri.org/">
//      <ORDERS_ADDResult>
//     		<RTN_ORDER_DOC>
//    			<RTN_RESULT>
//   				<PROCDATE>2020-11-16 11:40:25</PROCDATE>
//  				<PROCCNT>1</PROCCNT>
// 					<ERRCNT>0</ERRCNT>
//				</RTN_RESULT>
//				<ORDER>
//					<ECNO>461</ECNO><ODNO>L6261872894</ODNO>
//					<VENDORNO>20201116001</VENDORNO><ERRCODE>000</ERRCODE>
//					<ERRDESC><![CDATA[訂單接收成功，可寄件]]></ERRDESC>
//				</ORDER>
//			</RTN_ORDER_DOC>
//		</ORDERS_ADDResult>
//    </ORDERS_ADDResponse>
//  </soap:Body>
//</soap:Envelope>

type OrdersAddEnvelope struct {
	Body   OrdersAddBody  `xml:"Body"`
}

type OrdersAddBody struct {
	Body OrdersAddResponse `xml:"ORDERS_ADDResponse"`
}

type OrdersAddResponse struct {
	Body OrdersAddResult `xml:"ORDERS_ADDResult"`
}

type OrdersAddResult struct {
	//XMLName   xml.Name `xml:"ORDERS_ADDResult"`
	Body   OrdersAddResultRtnDoc `xml:"RTN_ORDER_DOC"`
}

type OrdersAddResultRtnDoc struct {
	//XMLName   xml.Name
	Return OrdersAddResultRtnResult `xml:"RTN_RESULT"`
	Order  OrdersAddResultOrder `xml:"ORDER"`
}

type OrdersAddResultRtnResult struct {
	XMLName   xml.Name
	ProcDate string   `xml:"PROCDATE"`
	ProcCnt  string   `xml:"PROCCNT"`
	ErrCnt   string   `xml:"ERRCNT"`
}

type OrdersAddResultOrder struct {
	XMLName   xml.Name
	EcNo     string   `xml:"ECNO"`
	OdNo     string   `xml:"ODNO"`
	VendorNO string   `xml:"VENDORNO"`
	ErrCode  string   `xml:"ERRCODE"`
	ErrDesc  string   `xml:"ERRDESC"`
}

func (receiver *OrdersAddEnvelope) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data,receiver)
	return
}
