package OKMart

import (
	"api/services/Enum"
	"api/services/entity"
	"api/services/util/log"
	m_sort "api/services/util/m-sort"
	"api/services/util/tools"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	client       http.Client
	ApiHost      string
	FtpHost      string
	FtpUsername  string
	FtpPassword  string
	KeyCode      string
	EcNo1        string
	EcNo2        string
	EcNo3        string
	EcNo4        string
	EcNo5        string
	F60StserCode string
	F60SerCode   string
	F10SerCode   string
	ConnectFtp      *ftp.ServerConn
	ConnectFtpError error
}

func (receiver *Client) GetClient() {
	config := viper.GetStringMapString(`MartOK`)
	receiver.ApiHost = config[`apihost`]
	receiver.EcNo1 = config[`ecno1`]
	receiver.EcNo2 = config[`ecno2`]
	receiver.EcNo3 = config[`ecno3`]
	receiver.EcNo4 = config[`ecno4`]
	receiver.EcNo5 = config[`ecno5`]
	receiver.F60StserCode = config[`f60stsercode`]
	receiver.F60SerCode = config[`f60sercode`]
	receiver.F10SerCode = config[`f10sercode`]
	receiver.FtpHost = config[`ftphost`]
	receiver.FtpUsername = config[`ftpusername`]
	receiver.FtpPassword = config[`ftppassword`]
	receiver.KeyCode = config[`keycode`]

	receiver.ConnectFtp, receiver.ConnectFtpError = newFtpAndLogin(receiver)
}

/**
OK Ftp 登入
*/
func newFtpAndLogin(receiver *Client) (c *ftp.ServerConn, err error) {
	c, err = ftp.Dial(receiver.FtpHost, ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(false))
	if err != nil {
		log.Error("OK Ftp 連線失敗 Error: [%]", err.Error())
		return c, fmt.Errorf("全家 Ftp 連線失敗")
	}

	err = c.Login(receiver.FtpUsername, receiver.FtpPassword)
	if err != nil {
		log.Error("OK Ftp 登入 Error: [%]", err.Error())
		return c, fmt.Errorf("OK Ftp 登入失敗")
	}

	return c, nil
}

/**
找尋資料夾底下檔案
*/
func (receiver *Client) FetchFolder(ecNoFolderName string) (fileNames []string, err error) {
	log.Info("OK  FetchFolder [%s]", ecNoFolderName)
	err = receiver.ConnectFtp.ChangeDir(ecNoFolderName)
	if err != nil {
		log.Error(ecNoFolderName + "OK FetchFolder ChangeDir WORK Fail")
		return fileNames, err
	}

	fileNames, err = receiver.ConnectFtp.NameList(".")

	if len(fileNames) == 0 {
		log.Error(ecNoFolderName + ` FetchFolder 無資料`)
		return fileNames, fmt.Errorf(ecNoFolderName + ` FetchFolder 無資料`)
	}

	fileNames = stringsRemoveString(fileNames, "./backup")

	log.Debug(`fileNames stringsRemoveString`, fileNames)
	return fileNames, err
}

/**
下載並轉碼吐出資料
*/
func (receiver *Client) RetrFile(ecNoFolderName, fileName string) ([]byte, error) {
	receiver.ConnectFtp.ChangeDir(ecNoFolderName)
	resp, err := receiver.ConnectFtp.Retr(fileName)

	path := fmt.Sprintf(`%s%s%s/`, viper.GetString(`Data.cvsPath`), `OK`, ecNoFolderName)
	log.Info(`RetrFile`, path)

	if err != nil {
		log.Error(ecNoFolderName + "載入失敗 [%s]", fileName)
		log.Error("Retr Error [%v]", err)
		return nil, fmt.Errorf(ecNoFolderName + "載入失敗 [%s]", fileName)
	}

	defer resp.Close()

	data, err := ioutil.ReadAll(resp)

	if err != nil {
		log.Error(ecNoFolderName + "ReadAll 失敗 [%s]", fileName)
		log.Error("ReadAll Error [%v]", err)
		return nil, err
	}

	convertData, err := tools.Big5ToUtf8ByByte(data)
	// 寫檔案
	tools.WriteFileByByte(path, fileName, convertData)

	return convertData, err
}

/**
處理後搬移位置
*/
func (receiver *Client) MoveFileToBackup(ecNoFolderName, fileName string) {
	err := receiver.ConnectFtp.ChangeDir(ecNoFolderName)
	if err != nil {
		log.Error(ecNoFolderName + ` MoveFileToBackup ChangeDir `, err.Error())
	}

	workF := ecNoFolderName + `/` + fileName
	destF := ecNoFolderName + `/backup/` + fileName
	err = receiver.ConnectFtp.Rename(workF, destF)
	log.Info(`OK Move`, workF, destF)
	if err != nil {
		log.Error("Rename [%s]", workF, destF)
		log.Error("Rename Error [%s]", err.Error())
		return
	}
}

func (receiver *Client) GetFlowType(ecNo string) string {

	if ecNo == receiver.EcNo5 || ecNo == receiver.EcNo3 {
		return `R` //逆向
	}

	return `N` // 順向
}

func (receiver *Client) OrderAdd(orderData entity.OrderData, sellerData entity.MemberData) (OrdersAddResult, error) {
	var amt, trType = `0`, `3` // 代收金額 | 取貨不付款

	if orderData.PayWay == Enum.CvsPay {
		amt = strconv.Itoa(int(orderData.TotalAmount))
		trType = `1`
	}

	request := `
	<root><VENDORCODE><KEYCODE><![CDATA[` + receiver.KeyCode + `]]></KEYCODE></VENDORCODE>
	<ORDER_DOC><ORDER>
			<ECNO>` + receiver.EcNo2 + `</ECNO>
			<STECNO>` + receiver.EcNo1 + `</STECNO>
			<ODNO></ODNO>
			<STNO>` + orderData.ReceiverAddress + `</STNO>
			<AMT>` + amt + `</AMT>
			<CUTKNM><![CDATA[` + orderData.ReceiverName + `]]></CUTKNM>
			<CUTKTL>` + orderData.ReceiverPhone + `</CUTKTL>
			<PRODNM>0</PRODNM>
			<ECWEB><![CDATA[` + sellerData.Username + `]]></ECWEB>
			<ECSERTEL>` + sellerData.Mphone + `</ECSERTEL>
			<REALAMT>` + strconv.Itoa(int(orderData.TotalAmount)) + `</REALAMT>
			<TRADETYPE>` + trType + `</TRADETYPE>
			<SERCODE>` + receiver.F60SerCode + `</SERCODE>
			<EDCNO>D13</EDCNO>
			<VENDOR><![CDATA[` + "Check'Ne" + `]]></VENDOR>
			<VENDORNO>` + orderData.OrderId + `</VENDORNO>
			<ORDERMODE>A</ORDERMODE>
			<PINCODE></PINCODE>
			<STAMT>0</STAMT>
			<STSERCODE>` + receiver.F60StserCode + `</STSERCODE>
	</ORDER></ORDER_DOC>
	<ORDERCOUNT><TOTALS>1</TOTALS></ORDERCOUNT></root> 
	`
	body := OrdersAddEnvelope{}
	result := OrdersAddResult{}
	finalData := privateSoapAndBase64("ORDERS_ADD", []byte(request))
	log.Debug("OrderAdd:", string(finalData))
	buf := bytes.NewBuffer(finalData)
	h := receiver.ApiHost + "/EC_C2C_WS/Service_EC.asmx"
	log.Debug("OrderAdd:", h)
	resp, err := receiver.client.Post(h, "text/xml; charset=utf-8", buf)
	if err != nil {
		return result, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	unescap := html.UnescapeString(string(data))
	log.Debug("OrderAdd:", unescap)
	err = body.DecodeXML([]byte(unescap))
	if err != nil {
		return result, err
	}
	result = body.Body.Body.Body
	return result, nil
}

func (receiver *Client) OrderPrintX(formData []string) ([]byte, error) {
	h := receiver.ApiHost + "/ECShippingOrders/Printer_B2C_batchPDF?FormData=" + strings.Join(formData, ",")
	log.Info(`OK OrderPrintX`, h)
	resp, err := receiver.client.Get(h)

	log.Info(`OK OrderPrintX resp`, resp)

	if err != nil {
		log.Error(`OK OrderPrintX 超商印單失敗`, err)
		return nil, fmt.Errorf("OK 超商印單失敗")
	}

	return ioutil.ReadAll(resp.Body)
}

func (receiver *Client) OrderPrint(pingCode, ecSerTel []string) ([]byte, error) {
	h := receiver.ApiHost + "/ECShippingOrders/Printer_B2C_batchPDF"
	argStr := []string{}
	for i, v := range pingCode {
		if len(ecSerTel) > i {
			arg := v + ":TOK:" + ecSerTel[i]
			argStr = append(argStr, arg)
		}
	}
	query := "?FormData=" + strings.Join(argStr, ",")
	resp, err := receiver.client.Get(h + query)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (receiver *Client) OrderReSend(ecNo, shipNo, nStoreId, rName, rPhone, price, needPay string) (OrdersResendResult, error) {
	amt, trType := "0", "3" // 代收金額
	if needPay == `1` {
		amt = price
		trType = "1"
	}
	request := `
	<root><VENDORCODE><KEYCODE><![CDATA[` + receiver.KeyCode + `]]></KEYCODE></VENDORCODE>
	<ORDER_DOC>
		<ORDER>
			<ECNO>` + ecNo + `</ECNO>
			<ODNO>` + shipNo + `</ODNO>
			<STNO>` + nStoreId + `</STNO>
			<AMT>` + amt + `</AMT>
			<CUTKNM><![CDATA[` + rName + `]]></CUTKNM>
			<CUTKTL>` + rPhone + `</CUTKTL>
			<PRODNM>0</PRODNM>
			<REALAMT>` + price + `</REALAMT>
			<TRADETYPE>` + trType + `</TRADETYPE>
			<SERCODE>` + receiver.F10SerCode + `</SERCODE>
			<EDCNO>D13</EDCNO>
		</ORDER>
	</ORDER_DOC>
	<ORDERCOUNT><TOTALS>1</TOTALS></ORDERCOUNT></root> 
	`

	body := OrdersResendEnvelope{}
	result := OrdersResendResult{}
	finalData := privateSoapAndBase64("ORDERS_RESEND", []byte(request))
	buf := bytes.NewBuffer(finalData)
	h := receiver.ApiHost + "/EC_C2C_WS/Service_EC.asmx"
	resp, err := receiver.client.Post(h, "text/xml; charset=utf-8", buf)
	if err != nil {
		return result, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	unescap := html.UnescapeString(string(data))
	err = body.DecodeXML([]byte(unescap))
	if err != nil {
		return result, err
	}
	result = body.Body.Body.Body
	return result, nil
}

func (receiver *Client) GetRawDocument(ecNo string, t string) (raw string, fileName string, e error) {
	folderName := "/" + ecNo + "/" + t
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return "", "", err
	}
	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return "", "", err
	}
	return string(utf8), fName, nil
}

func (receiver *Client) GetF01Document() (result F01Doc, e error) {
	folderName := "461/F01"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

/**
取得F01
*/
func (receiver *Client) GetF01Xml(data []byte) (_xml F01Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

/**
	取得F27
 */
func (receiver *Client) GetF27Xml(data []byte) (_xml F27Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

/**
	取得F25
*/
func (receiver *Client) GetF25Xml(data []byte) (_xml F25Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF84Xml(data []byte) (_xml F84Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF71Xml(data []byte) (_xml F71Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF63Xml(data []byte) (_xml F63Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF03Xml(data []byte) (_xml F03Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF67Xml(data []byte) (_xml F67Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF44Xml(data []byte) (_xml F44Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF64Xml(data []byte) (_xml F64Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF04Xml(data []byte) (_xml F04Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF17Xml(data []byte) (_xml F17Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF65Xml(data []byte) (_xml F65Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF05Xml(data []byte) (_xml F05Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}

func (receiver *Client) GetF07Xml(data []byte) (_xml F07Doc, err error) {
	if _xml.DecodeXML(data) != nil {
		return _xml, fmt.Errorf("解析失敗")
	}
	return _xml, nil
}


func (receiver *Client) GetF27Document(ecNo string) (result F27Doc, e error) {
	folderName := "/" + ecNo + "/F27"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF25Document(ecNo string) (result F25Doc, e error) {
	folderName := "/" + ecNo + "/F25"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF84Document(ecNo string) (result F84Doc, e error) {
	folderName := "/" + ecNo + "/F84"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF71Document(ecNo string) (result F71Doc, e error) {
	folderName := "/" + ecNo + "/F71"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF63Document(ecNo string) (result F63Doc, e error) {
	folderName := "/" + ecNo + "/F63"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF44Document(ecNo string) (result F44Doc, e error) {
	folderName := "/" + ecNo + "/F44"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF64Document(ecNo string) (result F64Doc, e error) {
	folderName := "/" + ecNo + "/F64"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF17Document(ecNo string) (result F17Doc, e error) {
	folderName := "/" + ecNo + "/F17"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF65Document(ecNo string) (result F65Doc, e error) {
	folderName := "/" + ecNo + "/F65"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF67Document(ecNo string) (result F67Doc, e error) {
	folderName := "/" + ecNo + "/F67"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF03Document(ecNo string) (result F03Doc, e error) {
	folderName := "/" + ecNo + "/F03"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF04Document(ecNo string) (result F04Doc, e error) {
	folderName := "/" + ecNo + "/F04"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF07Document(ecNo string) (result F07Doc, e error) {
	folderName := "/" + ecNo + "/F07"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetF05Document(ecNo string) (result F05Doc, e error) {
	folderName := "/" + ecNo + "/F05"
	data, fName, err := receiver.privateGetFile(folderName)
	if err != nil {
		return result, err
	}

	utf8, err := receiver.privateBig5ToUtf8(data)
	if err != nil {
		return result, err
	}

	err = result.DecodeXML(utf8)
	if err != nil {
		return result, err
	}

	err = receiver.privateMoveWorkToOk(folderName, fName)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (receiver *Client) GetAcctDocument() (result, fName string, e error) {
	//folderName := "/"+receiver.EcNo2+"/ACCT"
	//data, fName, err := receiver.privateGetFile(folderName)
	//if err != nil {
	//	return result, fName, err
	//}
	//
	//utf8,err := receiver.privateBig5ToUtf8(data)
	//if err != nil {
	//	return result, fName, err
	//}
	//
	//
	//err = receiver.privateMoveWorkToOk(folderName,fName)
	//if err != nil {
	//	return result, fName, err
	//}

	return result, fName, nil
}

// 取得資料夾內的排序最小的XML
func (receiver *Client) privateGetFile(folderName string) ([]byte, string, error) {
	logStr := ""
	f, err := receiver.newFtpAndLogin()
	if err != nil {
		return nil, "", err
	}
	defer f.Quit()

	err = f.ChangeDir(folderName)
	if err != nil {
		return nil, "", errors.New(logStr + " Error:" + err.Error())
	}

	entries, err := f.NameList(".")

	log.Debug(folderName + `OK GetFile [%v]`, entries)

	if err != nil {
		return nil, "", errors.New(logStr + " Error:" + err.Error())
	}

	entries = stringsRemoveString(entries, "./backup")

	log.Debug(folderName + `OK GetFile stringsRemoveString [%v]`, entries)

	oldName, err := m_sort.SortStringsAscAndGetFirst(entries)
	if err != nil {
		return nil, "", errors.New("File not found")
	}

	resp, err := f.Retr(oldName)
	if err != nil {
		return nil, "", err
	}

	data, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, "", err
	}

	return data, strings.TrimPrefix(oldName, "./"), nil
}

func (receiver *Client) privateBig5ToUtf8(data []byte) ([]byte, error) {
	big5ToUTF8 := traditionalchinese.Big5.NewDecoder()
	utf8, _, err := transform.Bytes(big5ToUTF8, data)
	return utf8, err
}

func (receiver *Client) newFtpAndLogin() (c *ftp.ServerConn, err error) {
	c, err = ftp.Dial(receiver.FtpHost, ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(false))
	if err != nil {
		return nil, err
	}

	err = c.Login(receiver.FtpUsername, receiver.FtpPassword)
	if err != nil {
		c.Quit()
		return nil, err
	}
	return c, nil
}

func (receiver *Client) privateMoveWorkToOk(docType, fileName string) (err error) {
	logStr := ""
	f, err := receiver.newFtpAndLogin()
	if err != nil {
		return
	}
	defer f.Logout()
	logStr += "Login[OK] "

	newFileName := strings.TrimPrefix(fileName, "./")

	workF := docType + "/"
	okF := workF + "backup/" + newFileName
	workFile := workF + newFileName

	err = f.Rename(workFile, okF)
	if err != nil {
		return err
	}
	return nil
}

func privateSoapAndBase64(opration string, data []byte) []byte {
	encodeData := base64.StdEncoding.EncodeToString(data)
	soapData := `<?xml version="1.0" encoding="utf-8"?>
		 <soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  		 <soap:Body><` + opration + ` xmlns="http://tempuri.org/">
		 <f>` + encodeData + `</f>
    	 </` + opration + `></soap:Body></soap:Envelope>
	`
	return []byte(soapData)
}

func stringsRemoveString(strs []string, str string) (final []string) {
	for _, v := range strs {
		if v == str {
			continue
		}
		final = append(final, strings.TrimPrefix(v, "./"))
	}
	return
}
