package curl

import (
	"api/services/VO/Response"
	"api/services/util/log"
	"api/services/util/tools"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var client *http.Client

func init() {
	client = &http.Client{}
}

func PostTransferXml(uri string, body string) ([]byte, error) {
	bodyEncode := url.QueryEscape(body)
	bodyBuffer := bytes.NewBuffer([]byte(bodyEncode))

	request, errorResponse := http.NewRequest("POST", uri, bodyBuffer)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Soapaction", "")

	//This is for test
	client := &http.Client{}
	resp, err := client.Do(request)
	if errorResponse != nil {
		log.Error("url error", errorResponse)
		return nil, err
	}
	//log.Info("Post tls", GetTLSVersion(resp.TLS))
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Debug("url body =>", string(body))
		return body, nil
	}
	return nil, err
}

func PostXml(uri string, body string) ([]byte, error) {
	bodyEncode := url.QueryEscape(body)
	bodyEncode = "xml=" + bodyEncode
	log.Debug("xml", bodyEncode)
	bodyBuffer := bytes.NewBuffer([]byte(bodyEncode))
	resp, err := client.Post(uri, "application/x-www-form-urlencoded", bodyBuffer)
	if err != nil {
		log.Error("post xml error", err)
		return nil, err
	}
	//log.Info("Post tls", GetTLSVersion(resp.TLS))
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Debug("post xml response =>", string(body))
		return body, nil
	}
	return nil, nil
}

func PostJson(uri string, body interface{}) ([]byte, error) {
	bodyEncode, _ := tools.JsonEncode(body)
	bodyBuffer := bytes.NewBuffer([]byte(bodyEncode))
	resp, err := client.Post(uri, "application/json", bodyBuffer)
	if err != nil {
		log.Error("post json error", err, bodyEncode, uri)
		return nil, err
	}
	//log.Info("Post tls", GetTLSVersion(resp.TLS))
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Debug("post json response =>", string(body))
		return body, nil
	}
	return nil, nil
}


func Post(uri string, body string) ([]byte, error) {

	bodyBuffer := bytes.NewBuffer([]byte(body))
	resp, err := client.Post(uri, "application/x-www-form-urlencoded", bodyBuffer)
	if err != nil {
		log.Error("post error", err)
		return nil, err
	}
	//log.Info("Post tls", GetTLSVersion(resp.TLS))
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		//log.Debug("post response =>", string(body))
		return body, nil
	}
	return nil, nil
}


func PostValues(uri string, body url.Values) ([]byte, error) {
	URL := fmt.Sprintf("%s?%s", uri, body.Encode())
	data := strings.NewReader(body.Encode())
	log.Debug("data => ", URL)
	resp, err := client.Post(URL, "", data)
	if err != nil {
		log.Error("post xml error", err)
		return nil, err
	}
	//log.Info("Post tls", GetTLSVersion(resp.TLS))
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Debug("post response =>", string(body))
		return body, nil
	}
	return nil, nil
}

func GetIDCheck(url string, header string) (Response.TWIDVerify, error) {

	var verify Response.TWIDVerify

	request, errorResponse := http.NewRequest("GET", url, nil)
	request.Header.Add("sris-consumerAdminId","00000000")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", header)

	//This is for test
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	inFunClient := &http.Client{Transport: tr}

	resp, err := inFunClient.Do(request)
	if err != nil {
		log.Error("url error", errorResponse)
		return verify, err
	}

	//log.Info("Post tls", GetTLSVersion(resp.TLS))
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &verify)
		log.Debug("url body =>", string(body))
		return verify, nil
	}
	return verify, nil
}

func Get(url string) ([]byte, error) {

	resp, err := client.Get(url)
	if err != nil {
		//log.Error("url error")
		return nil, err
	}
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Debug("url body =>", string(body))
		return body, nil
	}
	return nil, err
}

//func GetTLSVersion(tr *tls.ConnectionState) string {
//	log.Debug("sss", string(tr.Version))
//	if tr.Version > 0 {
//		switch tr.Version {
//		case tls.VersionSSL30:
//			return "SSL"
//		case tls.VersionTLS10:
//			return "TLS 1.0"
//		case tls.VersionTLS11:
//			return "TLS 1.1"
//		case tls.VersionTLS12:
//			return "TLS 1.2"
//		case tls.VersionTLS13:
//			return "TLS 1.3"
//		default:
//			return "Not"
//		}
//	}
//	return ""
//}
