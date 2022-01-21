package OKMart

import (
	"errors"
	"testing"
)

func TestF67Doc_DecodeXML(t *testing.T) {
	dataStr := `<F67DOC>
    <DOCHEAD>
        <DOCDATE>20201116</DOCDATE>
        <FROMPARTNERCODE>CVS</FROMPARTNERCODE>
        <TOPARTENERCODE>462</TOPARTENERCODE>
    </DOCHEAD>
    <F67CONTENT>
		  <RETM>T08</RETM>
		  <ECNO>516</ECNO>
		  <STNO>K001234</STNO>
		  <ODNO>12345678901</ODNO>
		  <RTDCDT>20180913144333</RTDCDT>
		  <FRTDCDT>20180902144333</FRTDCDT>
		  <VENDOR>EC 廠商</VENDOR>
		  <VENDORNO>TW1234567890</VENDORNO>
	</F67CONTENT>
    <F67CONTENT>
		  <RETM>T08</RETM><ECNO>516</ECNO><STNO>K001234</STNO><ODNO>12345678901</ODNO>
		  <RTDCDT>20180913144333</RTDCDT><FRTDCDT>20180902144333</FRTDCDT><VENDOR>EC 廠商</VENDOR>
		  <VENDORNO>TW1234567890</VENDORNO>
	</F67CONTENT>
	</F67DOC>`
	data := []byte(dataStr)

	doc := F67Doc{}
	err := doc.DecodeXML(data)
	if err != nil {
		t.Fatalf("DecodeXML() error = %v", err)
	}

	head := doc.Head
	if head.DocDate != "20201116" || head.ToPartenerCode != "462" || head.FromPartnerCode != "CVS" {
		t.Fatalf("DecodeXML() error = %v", errors.New("DocDate not match"))
	}

	if len(doc.Body) != 2 {
		t.Fatalf("DecodeXML() error = %v", errors.New("Content Length not match"))
	}
	content := doc.Body[0]
	if  content.ReturnCode != "T08" || content.EcNo != "516" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Barcode not match"))
	}

	if  content.RrStoreId != "K001234" || content.OrderNo != "12345678901" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

	if  content.UpDateTime != "20180913144333" || content.CheckDateTime != "20180902144333" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

	if  content.Vendor != "EC 廠商" || content.VendorNo != "TW1234567890" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}
}
