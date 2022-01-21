package OKMart

import (
	"errors"
	"testing"
)

func TestF44Doc_DecodeXML(t *testing.T) {
	dataStr := `<F44DOC>
    <DOCHEAD>
        <DOCDATE>20201116</DOCDATE>
        <FROMPARTNERCODE>CVS</FROMPARTNERCODE>
        <TOPARTENERCODE>462</TOPARTENERCODE>
    </DOCHEAD>
	<F44CONTENT>
		<ECNO>516</ECNO><ODNO>12345678901</ODNO><STNO>K001234</STNO><DCSTDT>20180913030035</DCSTDT><VENDOR>EC 廠商</VENDOR><VENDORNO>TW1234567890</VENDORNO>
	</F44CONTENT>
    <F44CONTENT>
		<ECNO>516</ECNO>
		<ODNO>12345678901</ODNO>
		<STNO>K001234</STNO>
		<DCSTDT>20180913030035</DCSTDT>
		<VENDOR>EC 廠商</VENDOR><VENDORNO>TW1234567890</VENDORNO>
	</F44CONTENT>
	</F44DOC>`
	data := []byte(dataStr)

	doc := F44Doc{}

	err := doc.DecodeXML(data)
	if err != nil {
		t.Fatalf("DecodeXML() error = %v", err)
	}

	if doc.Head.DocDate != "20201116" || doc.Head.ToPartenerCode != "462" || doc.Head.FromPartnerCode != "CVS" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Head not match"))
	}

	if len(doc.Body) != 2 {
		t.Fatalf("DecodeXML() error = %v", errors.New("Content Length not match"))
	}
	content := doc.Body[0]
	if  content.EcNo != "516" || content.OrderNo != "12345678901" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Barcode not match"))
	}

	if  content.RrStoreId != "K001234" || content.RrInDateTime != "20180913030035" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

	if  content.Vendor != "EC 廠商" || content.VendorNo != "TW1234567890" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}
}
