package OKMart

import (
	"errors"
	"testing"
)

func TestF63Doc_DecodeXML(t *testing.T) {
	dataStr := `<F63DOC>
	<DOCHEAD>
		<DOCDATE>20180913</DOCDATE>
		<FROMPARTNERCODE>CVS</FROMPARTNERCODE>
		<TOPARTNERCODE>516</TOPARTNERCODE>
	</DOCHEAD>
	<F63CONTENT>
		<ECNO>516</ECNO> <ODNO>12345678901</ODNO>
		<CNNO>TOK</CNNO> <CUNAME>簡大翔</CUNAME> <PRODTYPE>0</PRODTYPE>
		<PINCODE>12345678901</PINCODE>
		<RTDT>20180913163005</RTDT>
	</F63CONTENT>
  	<F63CONTENT>
		<ECNO>516</ECNO>
		<ODNO>12345678901</ODNO>
		<CNNO>TOK</CNNO> <CUNAME>簡大翔</CUNAME>
		<PRODTYPE>0</PRODTYPE>
		<PINCODE>12345678901</PINCODE>
		<RTDT>20180913163005</RTDT>
	</F63CONTENT>
	</F63DOC>`
	data := []byte(dataStr)

	doc := F63Doc{}
	err := doc.DecodeXML(data)
	if err != nil {
		t.Errorf("DecodeXML() error = %v", err)
	}

	head := doc.Head
	if head.DocDate != "20180913" || head.ToPartenerCode != "516" || head.FromPartnerCode != "CVS" {
		t.Fatalf("DecodeXML() error = %v", errors.New("DocDate not match"))
	}

	if len(doc.Body) != 2 {
		t.Fatalf("DecodeXML() error = %v", errors.New("Content Length not match"))
	}
	content := doc.Body[0]
	if  content.EcNo != "516" || content.OrderNo != "12345678901" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Barcode not match"))
	}

	if  content.CnNo != "TOK" || content.RrName != "簡大翔" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

	if  content.SendNo != "12345678901" || content.ProductType != "0" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

	if  content.UpDateTime != "20180913163005" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}
}
