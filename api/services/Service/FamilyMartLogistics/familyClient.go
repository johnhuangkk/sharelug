package FamilyMartLogistics

import (
	"api/services/VO/FamilyMart"
	"api/services/util/log"
	"api/services/util/tools"
	"api/services/util/unzip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
)

type Client struct {
	ApiKey          string
	PriKey          string
	Host            string
	FtpHost         string
	FtpUnzipPasswd  string
	FtpUsername     string
	FtpPassword     string
	client          http.Client
	ConnectFtp      *ftp.ServerConn
	ConnectFtpError error
}

func (receiver *Client) GetClient(reverseFlow bool) {
	var config map[string]string
	if reverseFlow {
		config = viper.GetStringMapString(`MartFamily911`)
	} else {
		config = viper.GetStringMapString(`MartFamily901`)
	}

	log.Info("GetClient: [%s]", config)

	receiver.Host = config[`apihost`]
	receiver.ApiKey = config[`apikey`]
	receiver.PriKey = config[`prikey`]
	receiver.FtpHost = config[`ftphost`]
	receiver.FtpUsername = config[`ftpusername`]
	receiver.FtpPassword = config[`ftppassword`]
	receiver.FtpUnzipPasswd = config[`ftpunzippasswd`]
	receiver.ConnectFtp, receiver.ConnectFtpError = newFtpAndLogin(receiver)
}

func (receiver *Client) OrderAdd(request FamilyMart.OrderAddRequest) (response FamilyMart.OrderAddResponse, err error) {
	api := "/C2COrderAdd/C2COrderAdd.ashx"

	ts, ru := receiver.GetTimeStamp()
	if !ru {
		return response, errors.New("OrderAdd.GetTimeStamp fail")
	}

	str, r := request.EncodeXML()
	if !r {
		return
	}
	log.Debug("OrderAdd.generateRequestBody:", receiver.ApiKey, receiver.PriKey, ts, str)
	body := generateRequestBody(receiver.ApiKey, receiver.PriKey, ts, str)
	buff := []byte(body)
	resp, err := receiver.client.Post(receiver.Host+api, "application/x-www-form-urlencoded", bytes.NewBuffer(buff))
	if err != nil {
		return response, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	log.Debug("OrderAdd:", string(data), err)
	if err != nil {
		return response, err
	}

	err = xml.Unmarshal(data, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (receiver *Client) OrderPrint(request OrdersPrintRequest) (imagePath string, result bool) {
	api := "/OrdersPrint/OrdersPrint.aspx"

	ts, ru := receiver.GetTimeStamp()
	if !ru {
		return
	}

	str, r := request.EncodeXML()
	if !r {
		return
	}

	body := generateRequestBody(receiver.ApiKey, receiver.PriKey, ts, str)
	buff := []byte(body)
	resp, err := receiver.client.Post(receiver.Host+api, "application/x-www-form-urlencoded", bytes.NewBuffer(buff))
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	content := string(data)
	if strings.Contains(content, "SendErrorPost.aspx") {
		return
	}

	content = strings.Split(content, `<img src="`)[1]
	content = strings.Split(content, `" width="`)[0]
	imagePath = receiver.Host + "/OrdersPrint/" + content

	return imagePath, true
}

func (receiver *Client) OrderSwitch(parentId, eshopId, shipNo, ecOrderNo, newStoreId, storeType string) (response OrderSwitchResponse, err error) {
	api := "/RCV_SWITCHSTORENOTIFY/RCV_SWITCHSTORENOTIFY.ashx"

	reqStr := `Data=<?xml version="1.0" encoding="UTF-8"?>
			<Doc>
				<ShipmentNos>
					<ParentId>` + parentId + `</ParentId>
					<EshopId>` + eshopId + `</EshopId>
					<OrderNo>` + shipNo + `</OrderNo>
					<EcOrderNo>` + ecOrderNo + `</EcOrderNo>
					<RcvStoreType>1</RcvStoreType>
					<RcvStoreId>` + newStoreId + `</RcvStoreId>
					<StoreType>` + storeType + `</StoreType>
				</ShipmentNos>
			</Doc>`
	log.Info("OrderSwitch reqStr [%s]", reqStr)
	buff := []byte(reqStr)
	resp, err := receiver.client.Post(receiver.Host+api, "application/x-www-form-urlencoded", bytes.NewBuffer(buff))
	if err != nil {
		return response, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	log.Debug("OrderSwitch:", string(data), err)
	if err != nil {
		return response, err
	}

	err = response.DecodeXML(data)
	return response, nil
}

func (receiver *Client) GetTimeStamp() (ts int64, result bool) {
	api := "/API_TIMESTAMP_QUERY/API_TIMESTAMP_QUERY.ashx"
	resp, err := receiver.client.Get(receiver.Host + api)
	if err != nil || resp.StatusCode != http.StatusOK {
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	tInt, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return
	}
	return tInt, true
}

/**
?????? I00 Xml
*/
func (receiver *Client) GetI00Xml(data []byte) (_xml FMLI00, err error) {

	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}
	return _xml, nil
}

/**
?????? R22 Xml
*/
func (receiver *Client) GetR22Xml(data []byte) (_xml FMLR22, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}
	return _xml, nil
}

/**
?????? R23 Xml
*/
func (receiver *Client) GetR23Xml(data []byte) (_xml FMLR23, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}
	return _xml, nil
}

/**
?????? R25 Xml
*/
func (receiver *Client) GetR25Xml(data []byte) (_xml FMLR25, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}
	return _xml, nil
}

/**
?????? R25 Xml
*/
func (receiver *Client) GetR27Xml(data []byte) (_xml FMLR27, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}
	return _xml, nil
}

func (receiver *Client) GetR04Xml(data []byte) (_xml FMLR04, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}
	return _xml, nil
}

func (receiver *Client) GetRS9Xml(data []byte) (_xml FMLRS9, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}
	return _xml, nil
}

func (receiver *Client) GetR28Xml(data []byte) (_xml FMLR28, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}
	return _xml, nil
}

func (receiver *Client) GetRS4Xml(data []byte) (_xml FMLRS4, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}

	return _xml, nil
}

func (receiver *Client) GetR29Xml(data []byte) (_xml FMLR29, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}

	return _xml, nil
}

func (receiver *Client) GetR96Xml(data []byte) (_xml FMLR96, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}

	return _xml, nil
}

func (receiver *Client) GetR08Xml(data []byte) (_xml FMLR08, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}

	return _xml, nil
}

func (receiver *Client) GetR89Xml(data []byte) (_xml FMLR89, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}

	return _xml, nil
}

func (receiver *Client) GetR98Xml(data []byte) (_xml FMLR98, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}

	return _xml, nil
}

func (receiver *Client) GetR99Xml(data []byte) (_xml FMLR99, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("????????????")
	}

	return _xml, nil
}

/**
???????????????????????????
*/
func (receiver *Client) FetchFolder(folderName string) (fileNames []string, err error) {
	log.Info("MartFamilyFetch [%s]", folderName)

	err = receiver.ConnectFtp.ChangeDir(`/SHARE28P/` + folderName + `/WORK`)
	if err != nil {
		log.Error(folderName + "FetchFolder ChangeDir WORK Fail")
		return fileNames, err
	}

	fileNames, err = receiver.ConnectFtp.NameList(".")

	if len(fileNames) == 0 {
		log.Error(folderName + ` FetchFolder ?????????`)
		return fileNames, fmt.Errorf(folderName + ` FetchFolder ?????????`)
	}

	return fileNames, err
}

/**
?????????????????????????????????
*/
func (receiver *Client) RetrFileAndUnzip(folderName, fileName string) ([]byte, error) {
	receiver.ConnectFtp.ChangeDir(`/SHARE28P/` + folderName + `/WORK`)
	resp, err := receiver.ConnectFtp.Retr(fileName)

	if err != nil {
		log.Error("Retr fileName [%v]", fileName)
		log.Error("Retr Error [%v]", err)
		return nil, fmt.Errorf("???????????? [%s]", folderName)
	}

	defer resp.Close()

	tempBuf, err := ioutil.ReadAll(resp)
	if err != nil {
		log.Error("ReadAll fileName [%v]", fileName)
		log.Error("ReadAll Error [%v]", err)
		return nil, err
	}

	data, err := unzip.UnzipDataWithPassword(tempBuf, receiver.FtpUnzipPasswd)
	path := fmt.Sprintf(`%s%s/%s/`, viper.GetString(`Data.cvsPath`), `Family`, folderName)
	// ?????????
	var re = regexp.MustCompile(`(.*).zip`)
	tools.WriteFileByByte(path, re.ReplaceAllString(fileName, `$1`), data)

	if err != nil {
		log.Error("UnzipDataWithPassword fileName [%s]", fileName)
		log.Error("UnzipDataWithPassword Error [%s]", err.Error())
		return nil, fmt.Errorf("fileName Unzip Fail")
	}

	return data, nil
}

/**
?????????????????????
*/
func (receiver *Client) MoveFileToDest(folderName, fileName, dest string) {
	var path = `/SHARE28P/` + folderName
	receiver.ConnectFtp.ChangeDir(path)

	workF := path + `/WORK/` + fileName
	destF := path + `/` + dest + `/` + fileName
	err := receiver.ConnectFtp.Rename(workF, destF)
	//
	if err != nil {
		log.Error("Rename [%s]", workF, destF)
		log.Error("Rename Error [%s]", err.Error())
		return
	}
}

/**
Family Ftp ??????
*/
func newFtpAndLogin(receiver *Client) (c *ftp.ServerConn, err error) {
	c, err = ftp.Dial(receiver.FtpHost, ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(false))
	if err != nil {
		log.Error("?????? Ftp ???????????? Error: [%]", err.Error())
		return c, fmt.Errorf("?????? Ftp ????????????")
	}

	err = c.Login(receiver.FtpUsername, receiver.FtpPassword)
	if err != nil {
		log.Error("?????? Ftp ?????? Error: [%]", err.Error())
		return c, fmt.Errorf("?????? Ftp ????????????")
	}

	return c, nil
}
