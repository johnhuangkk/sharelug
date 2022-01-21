package OKMart

import (
	"errors"
	"testing"
)

func TestF64Doc_DecodeXML(t *testing.T) {
	dataStr := `<F64DOC>
    <DOCHEAD>
        <DOCDATE>20201116</DOCDATE>
        <FROMPARTNERCODE>CVS</FROMPARTNERCODE>
        <TOPARTENERCODE>462</TOPARTENERCODE>
    </DOCHEAD>
    <F64CONTENT>
		<ECNO>516</ECNO>
		<STNO>K001234</STNO>
		<ODNO>12345678901</ODNO>
		<DCSTDT>20180913030035</DCSTDT>
		<VENDOR>EC 廠商</VENDOR>
		<VENDORNO>TW1234567890</VENDORNO>
	</F64CONTENT>
    <F64CONTENT>
		<ECNO>516</ECNO>
		<STNO>K001234</STNO>
		<ODNO>12345678901</ODNO>
		<DCSTDT>20180913030035</DCSTDT>
		<VENDOR>EC 廠商</VENDOR>
		<VENDORNO>TW1234567890</VENDORNO>
	</F64CONTENT>
	</F64DOC>`
	data := []byte(dataStr)

	doc := F64Doc{}

	err := doc.DecodeXML(data)
	if err != nil {
		t.Fatalf("DecodeXML() error = %v", err)
	}

	if doc.Head.DocDate != "20201116" {
		t.Fatalf("DecodeXML() error = %v", errors.New("DocDate not match"))
	}

	if doc.Head.ToPartenerCode != "462" {
		t.Fatalf("DecodeXML() error = %v", errors.New("ToPartenerCode not match"))
	}

	if doc.Head.FromPartnerCode != "CVS" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

	if len(doc.Body) != 2 {
		t.Fatalf("DecodeXML() error = %v", errors.New("Content Length not match"))
	}
	content := doc.Body[0]
	if  content.EcNo != "516" || content.RrStoreId != "K001234" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Barcode not match"))
	}

	if  content.OrderNo != "12345678901" || content.RrInDateTime != "20180913030035" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

	if  content.Vendor != "EC 廠商" || content.VendorNo != "TW1234567890" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}
}
