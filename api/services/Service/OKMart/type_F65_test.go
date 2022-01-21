package OKMart

import (
	"errors"
	"testing"
)

func TestF65Doc_DecodeXML(t *testing.T) {
	dataStr := `<F65DOC>
    <DOCHEAD>
        <DOCDATE>20201116</DOCDATE>
        <FROMPARTNERCODE>CVS</FROMPARTNERCODE>
        <TOPARTENERCODE>462</TOPARTENERCODE>
    </DOCHEAD>
    <F65CONTENT>
		<BC1>516725992</BC1><BC2>3916921910002084</BC2>
		<STNO>K001234</STNO>
		<RTDT>20180914140036</RTDT>
		<TKDT>20180914140036</TKDT><VENDOR>EC 廠商</VENDOR><VENDORNO>TW1234567890</VENDORNO></F65CONTENT>
    <F65CONTENT><BC1>516725992</BC1><BC2>3916921910002084</BC2><STNO>>K001234</STNO><RTDT>20180914140036</RTDT><TKDT>20180914140036</TKDT><VENDOR>EC 廠商</VENDOR><VENDORNO>TW1234567890</VENDORNO></F65CONTENT>
	</F65DOC>`
	data := []byte(dataStr)

	doc := F65Doc{}
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
	if  content.BarCode1 != "516725992" || content.BarCode2 != "3916921910002084" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Barcode not match"))
	}

	if  content.RrStoreId != "K001234" || content.RrPickDateTime != "20180914140036" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

	if  content.CheckDateTime != "20180914140036" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

	if  content.Vendor != "EC 廠商" || content.VendorNo != "TW1234567890" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}
}
