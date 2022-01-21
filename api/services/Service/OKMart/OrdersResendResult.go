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
//				<ERR_ORDER>
//					<ECNO>417</ECNO>
//					<ODNO>12345678901</ODNO>
//					<STNO>K001234</STNO>
//					<AMT>500</AMT>
//					<CUTKNM><![CDATA[簡大翔]]></CUTKNM>
//					<CUTKTL>123</CUTKTL>
//					<PRODNM>0</PRODNM>
//					<REALAMT>500</REALAMT>
//					<TRADETYPE>1</TRADETYPE>
//					<SERCODE>993</SERCODE>
//					<EDCNO>D13</EDCNO>
//					<ERRCODE>E06</ERRCODE>
//					<ERRDESC><![CDATA[踢退無店舖編號]]</ERRDESC>
//				</ERR_ORDER>
//			</RTN_ORDER_DOC>
//		</ORDERS_ADDResult>
//    </ORDERS_ADDResponse>
//  </soap:Body>
//</soap:Envelope>

type OrdersResendEnvelope struct {
	Body OrdersResendBody `xml:"Body"`
}

type OrdersResendBody struct {
	Body OrdersResendResponse `xml:"ORDERS_RESENDResponse"`
}

type OrdersResendResponse struct {
	Body OrdersResendResult `xml:"ORDERS_RESENDResult"`
}

type OrdersResendResult struct {
	//XMLName   xml.Name `xml:"ORDERS_ADDResult"`
	Body OrdersResendResultRtnDoc `xml:"RTN_ORDER_DOC"`
}

type OrdersResendResultRtnDoc struct {
	//XMLName   xml.Name
	Return OrdersResendResultRtnResult  `xml:"RTN_RESULT"`
	Order  OrdersResendResultErrorOrder `xml:"ERR_ORDER"`
}

type OrdersResendResultRtnResult struct {
	XMLName  xml.Name
	ProcDate string `xml:"PROCDATE"`
	ProcCnt  string `xml:"PROCCNT"`
	ErrCnt   string `xml:"ERRCNT"`
}

type OrdersResendResultErrorOrder struct {
	XMLName   xml.Name
	EcNo      string `xml:"ECNO"`
	OdNo      string `xml:"ODNO"`
	StNo      string `xml:"STNO"`
	AMT       string `xml:"AMT"`
	RName     string `xml:"CUTKNM"`
	RPhone    string `xml:"CUTKTL"`
	PRODNM    string `xml:"PRODNM"`
	Price     string `xml:"REALAMT"`
	TradeType string `xml:"TRADETYPE"`
	SerCode   string `xml:"SERCODE"`
	EDcNo     string `xml:"EDCNO"`
	ErrorCode string `xml:"ERRCODE"`
	ErrorDesc string `xml:"ERRDESC"`
}

func (receiver *OrdersResendEnvelope) DecodeXML(data []byte) (err error) {
	err = xml.Unmarshal(data, receiver)
	return
}
