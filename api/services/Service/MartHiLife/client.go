package MartHiLife

import (
	"api/services/util/SetupSFtp"
	"api/services/util/log"
	"api/services/util/tools"
	"api/services/util/unzip"
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Client struct {
	FtpHost         string
	FtpUnzipPasswd  string
	FtpUsername     string
	FtpPassword     string
	ApiHost1        string
	ApiHost2        string
	ConnectFtp      *sftp.Client
	ConnectFtpError error
}

func (receiver *Client) GetClient() *Client {
	config := viper.GetStringMapString(`MartHiLife`)

	receiver.ApiHost1 = config[`apihost1`]
	receiver.ApiHost2 = config[`apihost2`]
	receiver.FtpHost = config[`ftphost`]
	receiver.FtpPassword = config[`ftppassword`]
	receiver.FtpUsername = config[`ftpusername`]
	receiver.FtpUnzipPasswd = config[`ftpunzippasswd`]
	receiver.ConnectFtp, receiver.ConnectFtpError = SetupSFtp.SFTPConnect(receiver.FtpHost, receiver.FtpUsername, receiver.FtpPassword)

	return receiver
}

func (receiver *Client) OrderAdd(add OrderAddRequest) (response OrderAddResp, err error) {
	path := "/ecapi/v1/ec_orders_platform_P.aspx"
	client := http.Client{}

	buf := bytes.NewBuffer(add.EncodeXML())
	resp, err := client.Post(receiver.ApiHost1+path, "application/x-www-form-urlencoded", buf)
	if err != nil {
		return response, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	err = response.DecodeXML(data)

	return response, err
}

func (receiver *Client) OrderPrint(shipNo string) (data []byte, err error) {
	path := "/ecapi/v1/ec_ordersprn_C2C.aspx"
	client := http.Client{}

	str := `ParentId=124&EshopId=901&OrderNo=` + shipNo

	buf := bytes.NewBuffer([]byte(str))
	resp, err := client.Post(receiver.ApiHost2+path, "application/x-www-form-urlencoded", buf)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (receiver *Client) OrderSwitchStore(req OrderSwitchRequest) (response OrderSwitchResponse, err error) {
	hKey, hIv := viper.GetString("MartHiLife.HashKey"), viper.GetString("MartHiLife.HashIV")
	path := "/ecapi/v1/ec_orders_rcvswstore.aspx"
	chkMac := GenerateCheckSum3(req.ParentId, req.EshopId, req.EcDcNo, req.EcCvs, req.ShipNo, hKey, hIv)
	a := `Data=<?xml version="1.0" encoding="utf-8"?>
	<Doc><ShipmentNos>
		<ParentId>` + req.ParentId + `</ParentId>
		<EshopId>` + req.EshopId + `</EshopId>
		<EcDcNo>` + req.EcDcNo + `</EcDcNo>
		<EcCvs>` + req.EcCvs + `</EcCvs>
		<OrderNo>` + req.ShipNo + `</OrderNo>
		<EcOrderNo>` + req.EcOrderNo + `</EcOrderNo>
		<RcvStoreId>` + req.RcvStoreId + `</RcvStoreId>
		<StoreType>` + req.StoreType + `</StoreType>
		<ChkMac>` + chkMac + `</ChkMac>
		</ShipmentNos>
	</Doc>`
	client := http.Client{}

	buf := bytes.NewBuffer([]byte(a))
	resp, err := client.Post(receiver.ApiHost1+path, "application/x-www-form-urlencoded", buf)

	if err != nil {
		log.Error("client.Post Error [%v]", err)
		return response, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error("ioutil.ReadAll Error [%v]", err)
		return response, err
	}

	err = response.DecodeXML(data)
	log.Error("HiLife OrderSwitchStore response : [%v]", response)

	return response, err
}


/**
找尋資料夾底下檔案
*/
func (receiver *Client) FetchFolder(folderName, fileExt string) (fileNames []string, err error) {
	log.Info("MartHiLife Fetch Folder [%s]", folderName)
	files, err := receiver.ConnectFtp.ReadDir(`/` + folderName + `/WORK`)
	if err != nil {
		return fileNames, err
	}
	for _, v := range files {
		name := v.Name()
		if strings.HasPrefix(name, folderName) && strings.HasSuffix(name, fileExt) {
			fileNames = append(fileNames, name)
		}
	}

	if len(fileNames) == 0 {
		log.Info("Folder [%s] 無資料", folderName)
		return nil, fmt.Errorf("[%s] 資料夾無資料 ", folderName)
	}

	return fileNames, err
}

/**
下載及解壓縮並吐出資料
*/
func (receiver *Client) RetrFileAndUnzip(folderName, fileName string) ([]byte, error){
	resp, err := receiver.ConnectFtp.Open(`/` + folderName + `/WORK/` + fileName)

	if err != nil {
		log.Error("Open fileName [%s] Fail", fileName)
		log.Error("Open Error [%v]", err)
		return nil, fmt.Errorf("載入失敗 [%s]", folderName)
	}
	defer resp.Close()


	tempBuf, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	data, err := unzip.UnzipDataWithPassword(tempBuf, receiver.FtpUnzipPasswd)
	path := fmt.Sprintf(`%s%s/%s/`, viper.GetString(`Data.cvsPath`), `HiLife`, folderName)
	// 寫檔案
	var re = regexp.MustCompile(`(.*).zip`)
	tools.WriteFileByByte(path, re.ReplaceAllString(fileName, `$1`), data)

	if err != nil {
		log.Error("UnzipDataWithPassword fileName [%s] Fail", fileName)
		log.Error("UnzipDataWithPassword [%s]", err.Error())
		return nil, fmt.Errorf("fileName Unzip Fail")
	}

	return data ,nil
}

// 檔案搬移位置
func (receiver *Client) MoveFileToDest(folderName, fileName, dest string) {

	workF := fmt.Sprintf(`/%s/WORK/%s`, folderName, fileName)
	destF := fmt.Sprintf(`/%s/%s/%s`, folderName, dest, fileName)
	err := receiver.ConnectFtp.Rename(workF, destF)
	//
	if err != nil {
		log.Fatal(err.Error())
		log.Error("Rename  [%s]", workF, destF)
		return
	}
}

func (receiver *Client) GetR00Xml(data []byte) (_xml R00Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetR27Xml(data []byte) (_xml R27Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetR22Xml(data []byte) (_xml R22Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetR04Xml(data []byte) (_xml R04Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetR28Xml(data []byte) (_xml R28Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetRS4Xml(data []byte) (_xml RS4Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetR29Xml(data []byte) (_xml R29Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetR96Xml(data []byte) (_xml R96Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetR08Xml(data []byte) (_xml R08Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetRS9Xml(data []byte) (_xml RS9Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetR98Xml(data []byte) (_xml R98Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetR99Xml(data []byte) (_xml R99Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetR89Xml(data []byte) (_xml R89Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}