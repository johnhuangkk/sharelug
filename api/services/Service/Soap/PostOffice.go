package Soap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

const (
	uri   string = "http://ctewsp.post.gov.tw/PSTPT.WebService/PSTPT_WebService.asmx?op=TransferForAddressString"
	soap  string = "http://www.w3.org/2003/05/soap-envelope"
	xsi   string = "http://www.w3.org/2001/XMLSchema-instance"
	xsd   string = "http://www.w3.org/2001/XMLSchema"
	xmlns string = "http://tempuri.org/"
)

type soap12RQ struct {
	XMLName   xml.Name `xml:"soap12:Envelope"`
	XMLNsSoap string   `xml:"xmlns:soap12,attr"`
	XMLNsXSI  string   `xml:"xmlns:xsi,attr"`
	XMLNsXSD  string   `xml:"xmlns:xsd,attr"`
	Body      soap12Body
}

type RqData struct {
	XMLName  xml.Name `xml:"TransferForAddressString"`
	XMLNs    string   `xml:"xmlns,attr"`
	CAddress string   `xml:"cAddress"`
}

type soap12Body struct {
	XMLName xml.Name `xml:"soap12:Body"`
	Payload RqData
}

type soap12RS struct {
	XMLName   xml.Name `xml:"Envelope"`
	XMLNsSoap string   `xml:"xmlns:soap,attr"`
	XMLNsXSI  string   `xml:"xmlns:xsi,attr"`
	XMLNsXSD  string   `xml:"xmlns:xsd,attr"`
	Body      soap12RsBody
}

type RsData struct {
	XMLName                        xml.Name `xml:"TransferForAddressStringResponse"`
	XMLNs                          string   `xml:"xmlns,attr"`
	TransferForAddressStringResult string
}

type soap12RsBody struct {
	XMLName xml.Name `xml:"Body"`
	Payload RsData
}

func TranslateAddress(addr string) (string, error) {

	v := soap12RQ{
		XMLNsSoap: soap,
		XMLNsXSD:  xsd,
		XMLNsXSI:  xsi,
		Body: soap12Body{
			Payload: RqData{
				XMLNs:    xmlns,
				CAddress: addr,
			},
		},
	}

	payload, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return "", err
	}
	fmt.Printf("%q", dump)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(bodyBytes))
	res := soap12RS{}
	err = xml.Unmarshal(bodyBytes, &res)

	fmt.Println(res)
	fmt.Println(res.Body.Payload.TransferForAddressStringResult)
	defer response.Body.Close()
	return res.Body.Payload.TransferForAddressStringResult, nil
}
