package Soap

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
)

type soapRQ struct {
	XMLName   xml.Name `xml:"S:Envelope"`
	XMLNsSoap string   `xml:"xmlns:S,attr"`
	Body      soapBody
}

type soapBody struct {
	XMLName 	xml.Name	`xml:"S:Body"`
	GetData 	getData
}

type getData struct {
	XMLName	xml.Name `xml:"getData"`
	XMLNsSoap string `xml:"xmlns,attr"`
	XString xstring
}

type xstring struct {
	XMLName	xml.Name `xml:"xstring"`
	CData cdata `xml:"![CDATA["`
}
type cdata struct {
	SMX interface{}
}


func Call(ws string, payloadInterface interface{}) (string, error) {

	v := soapRQ{
		XMLNsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		Body: soapBody{
			GetData: getData{
				XMLNsSoap: "http://ws.kgibank.thesys.com",
				XString: xstring{
					CData: cdata{
						SMX: payloadInterface,
					},
				},
			},
		},
	}
	payload, err := xml.MarshalIndent(v, "", "  ")
	payload = []byte(xml.Header + string(payload))

	payload = bytes.Replace(payload, []byte("<![CDATA[>"), []byte("<![CDATA["), -1)
	payload = bytes.Replace(payload, []byte("</![CDATA[>"), []byte("]]>"), -1)

	//log.Debug("XML", string(payload))
	client := http.Client{}

	req, err := http.NewRequest("POST", ws, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("SOAPAction", "")
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	decodedValue := html.UnescapeString(string(bodyBytes))
	defer response.Body.Close()
	return decodedValue, nil
}
