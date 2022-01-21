package OKMart

import (
	"errors"
	"testing"
)

func TestF03Doc_DecodeXML(t *testing.T) {
	dataStr := `<F03DOC>
	<DOCHEAD>
		<DOCDATE>20180915</DOCDATE><FROMPARTNERCODE>CVS</FROMPARTNERCODE><TOPARTNERCODE>417</TOPARTNERCODE>
	</DOCHEAD>
	<F03CONTENT>
		<ECNO>417</ECNO><ODNO>12345678901</ODNO><CUNAME>簡大翔</CUNAME>
		<PRODTYPE>0</PRODTYPE><PINCODE>12345678901</PINCODE>
		<RTDT>20180915172249</RTDT>
	</F03CONTENT>
	<F03CONTENT>
		<ECNO>417</ECNO><ODNO>12345678901</ODNO><CUNAME>簡大翔</CUNAME>
		<PRODTYPE>0</PRODTYPE><PINCODE>12345678901</PINCODE>
		<RTDT>20180915172249</RTDT>
	</F03CONTENT>
	</F03DOC>`
	data := []byte(dataStr)

	doc := F03Doc{}
	err := doc.DecodeXML(data)
	if err != nil {
		t.Fatalf("DecodeXML() error = %v", err)
	}

	head := doc.Head
	if head.DocDate != "20180915" || head.ToPartenerCode != "417" || head.FromPartnerCode != "CVS" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Head not match"))
	}

	if len(doc.Body) != 2 {
		t.Fatalf("DecodeXML() error = %v", errors.New("Content Length not match"))
	}
	content := doc.Body[0]
	if  content.EcNo != "417" || content.OrderNo != "12345678901" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}

	if  content.RrName != "簡大翔" || content.ProductType != "0" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}

	if  content.SendNo != "12345678901" || content.UpDateTime != "20180915172249" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Body not match"))
	}
}
