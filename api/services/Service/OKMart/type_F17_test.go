package OKMart

import (
	"errors"
	"testing"
)

func TestF17Doc_DecodeXML(t *testing.T) {
	dataStr := `<F17DOC>
    <DOCHEAD>
        <DOCDATE>20201116</DOCDATE>
        <FROMPARTNERCODE>CVS</FROMPARTNERCODE>
        <TOPARTENERCODE>462</TOPARTENERCODE>
    </DOCHEAD>
	  <F17CONTENT>
	  	<BC1>516725992</BC1>
	  	<BC2>3916921910002084</BC2>
	  	<STNO>K001234</STNO>
	  	<RTDT>20180914140036</RTDT>
	  	<PINCODE></PINCODE>
	  </F17CONTENT>
	  <F17CONTENT>
	  	<BC1>516725992</BC1>
	  	<BC2>3916921910002084</BC2>
	  	<STNO>K001234</STNO>
	  	<RTDT>20180914140036</RTDT>
	  	<PINCODE></PINCODE>
	  </F17CONTENT>
	</F17DOC>`
	data := []byte(dataStr)

	doc := F17Doc{}

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
	if  content.BarCode1 != "516725992" || content.BarCode2 != "3916921910002084" {
		t.Fatalf("DecodeXML() error = %v", errors.New("Barcode not match"))
	}

	if  content.RrStoreId != "K001234" || content.RrPickDateTime != "20180914140036" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

	if  content.SendNo != "" {
		t.Fatalf("DecodeXML() error = %v", errors.New("FromPartnerCode not match"))
	}

}
