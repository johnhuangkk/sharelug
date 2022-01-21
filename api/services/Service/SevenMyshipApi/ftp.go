package SevenMyshipApi

import (
	"api/services/Enum"
	"api/services/Service/CvsShipping"
	"api/services/dao/Cvs"
	sevenmyshipdao "api/services/dao/SevenMyshipDao"
	"api/services/database"
	"api/services/entity"
	"api/services/util/SetupSFtp"
	"api/services/util/log"
	"bufio"
	"bytes"
	"encoding/xml"

	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pkg/sftp"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type CESPXml struct {
	Xml     xml.Name      `xml:"SPDoc"`
	DocHead CESPDocHeader `xml:"DocHead"`
	SP      []SP          `xml:"DocContent>SP"`
}
type CESPDocHeader struct {
	DocNo           string
	DocDate         string
	FromPartnerCode string `xml:"From>FromPartnerCode"`
	ToPartnerCode   string `xml:"To>ToPartnerCode"`
	DocCount        string
}

type SP struct {
	Xml         xml.Name `xml:"SP"`
	ParentId    string
	TotalCount  string
	TotalAmount string
	Detail      []SPDetail `xml:"SPDetail"`
}

type SPDetail struct {
	ParentId      string
	EshopId       string
	EC3GParentId  string
	EC3GEshopId   string
	PaymentNo     string
	DCStoreStatus string
	StoreId       string
	SPDate        string
	SPAmount      string
	ServiceType   string
	SPNo          string
}
type CEDRXml struct {
	Xml        xml.Name      `xml:"DCReceiveDoc"`
	DocHead    CEINDocHeader `xml:"DocHead"`
	DocContent []DCReturn    `xml:"DocContent>DCReturn"`
}
type EDRDocHeader struct {
	DocNo           string
	DocDate         string
	FromPartnerCode string `xml:"From>FromPartnerCode"`
	ToPartnerCode   string `xml:"To>ToPartnerCode"`
	DocCount        string
}
type DCReturn struct {
	EshopId        string
	EC3GParentId   string
	EC3GEshopId    string
	PaymentNo      string
	ShipmentNo     string
	ShipmentAmount string
	DCReturnDate   string
	DCReturnName   string
	DCReturnCode   string
}
type CERTXml struct {
	Xml        xml.Name         `xml:"DCReturnAdviceDoc"`
	DocHead    CERTDocHeader    `xml:"DocHead"`
	DocContent []DCReturnAdvice `xml:"DocContent>DCReturnAdvice"`
}
type CERTDocHeader struct {
	DocNo           string
	DocDate         string
	FromPartnerCode string `xml:"From>FromPartnerCode"`
	ToPartnerCode   string `xml:"To>ToPartnerCode"`
	DocCount        string
}
type DCReturnAdvice struct {
	EshopId             string
	EC3GParentId        string
	EC3GEshopId         string
	PaymentNo           string
	ShipmentNo          string
	ShipmentAmount      string
	ReturnType          string
	DCPlannedReturnDate string
}
type CEINXml struct {
	Xml        xml.Name      `xml:"DCReceiveDoc"`
	DocHead    CEINDocHeader `xml:"DocHead"`
	DocContent []DCReceive   `xml:"DocContent>DCReceive"`
}
type CEINDocHeader struct {
	DocNo           string
	DocDate         string
	FromPartnerCode string `xml:"From>FromPartnerCode"`
	ToPartnerCode   string `xml:"To>ToPartnerCode"`
	DocCount        string
}
type DCReceive struct {
	EshopId         string
	EC3GParentId    string
	EC3GEshopId     string
	PaymentNo       string
	ShipmentNo      string
	DCReceiveDate   string
	DCStoreStatus   string
	DCReceiveStatus string
	DCRecName       string
	DCStoreDate     string
}
type CPPSXml struct {
	Xml        xml.Name      `xml:"PPSDoc"`
	DocHead    CPPSDocHeader `xml:"DocHead"`
	DocContent []PPS         `xml:"DocContent>PPS"`
}
type CPPSDocHeader struct {
	DocNo           string
	DocDate         string
	FromPartnerCode string `xml:"From>FromPartnerCode"`
	ToPartnerCode   string `xml:"To>ToPartnerCode"`
	DocCount        string
}
type PPS struct {
	ParentId    string
	EshopId     string
	EC3GEshopId string
	PaymentNo   string
	ShipmentNo  string
	StoreId     string
	StoreDate   string
	StoreTime   string
	StoreType   string
	TelNo       string
}
type SevenShopData struct {
	ShopID   string
	ShopName string
	Country  string
	District string
	Address  string
}
type SevenShops struct {
	Shops []SevenShopData
}
type BoolTest struct {
	Test     string
	Bbb      bool
	Country  string
	District string
	Address  string
}

var flag bool

type Client struct {
	FtpHost         string
	FtpUsername     string
	FtpPassword     string
	ConnectFtp      *sftp.Client
	ConnectFtpError error
}

var folderPath = map[string]string{
	"Shop":                   "/CSTD", //每日營業店鋪地址資料
	"ReturnAccept":           "/CEDR", // 貨品已退到物流中心 逆
	"ReceiveAccept":          "/CEIN", // 貨品已交寄到物流中心 正
	"ReturnRequest":          "/CERT", // 買家未取預計退回通知
	"SellerSendStore":        "/OL",   // 賣家到店寄件紀錄
	"StoreChange":            "/CCS",  //寄/退店必轉通知
	"ArrivedShop":            "/CPPS", //商品到店
	"WeeklyOrderReport":      "/ACC",  //週結帳款
	"DailyShopOperateRecord": "/CESP", // 店鋪商品代收
	"LostAmount":             "/ACTR", //判賠檔案

}

var olSection = map[string]string{
	"1": "Start",
	"2": "Content",
	"3": "End",
}

func (receiver *Client) getClient() *Client {
	config := viper.GetStringMapString(`MYSHIP`)

	receiver.FtpHost = config[`sftphost`]
	receiver.FtpPassword = config[`sftppassword`]
	receiver.FtpUsername = config[`sftpaccount`]
	receiver.ConnectFtp, receiver.ConnectFtpError = SetupSFtp.SFTPConnect(receiver.FtpHost, receiver.FtpUsername, receiver.FtpPassword)

	return receiver
}
func GetShopAddress(country string, district string) ([]entity.SevenMyshipShopData, error) {
	address, err := sevenmyshipdao.FindShopsAddress(country, district)
	if err != nil {
		return address, err
	}
	return address, nil
}

func FetchPackageSendByOl() {
	now := time.Now()
	client := new(Client).getClient().ConnectFtp

	defer client.Close()

	files, err := client.ReadDir(folderPath["SellerSendStore"])
	if err != nil {
		log.Error("Seven Ftp Ol read error", err.Error())
		return
	}
	log.Info("sftp connect", "7-11 OL每日出貨交寄更新開始", now.Format(`2006-01-02 15:04:05`))
	tempDirPath, err := ioutil.TempDir("", "ol")
	if err != nil {
		log.Error("Seven Ftp Ol tmep dir create error", err.Error())
		return
	}

	i := 0
	for _, file := range files {
		if i > 5 {
			break
		}
		if !file.IsDir() {

			t, _ := time.Parse(`20060102`, strings.SplitAfter(file.Name(), ".")[1])

			if t.Before(now) {
				data, err := client.Open(folderPath["SellerSendStore"] + "/" + file.Name())

				if err != nil {
					log.Error("Seven Ftp OL open error"+file.Name(), err.Error())
					continue
				}
				tempfile, err := ioutil.TempFile(tempDirPath, "ol-*-"+file.Name()+".ol")
				if err != nil {
					log.Error("Seven Ftp OL  open temp error"+tempfile.Name(), err.Error())
					continue
				}
				data.WriteTo(tempfile)

				err = client.Rename(folderPath["SellerSendStore"]+"/"+file.Name(), folderPath["SellerSendStore"]+"/done/"+file.Name())
				if err != nil {
					log.Error("Seven Ftp OL open temp error"+tempfile.Name(), err.Error())
					continue
				}
			}

		}
		i++
	}
	client.Close()
	files, err = ioutil.ReadDir(tempDirPath)
	if err != nil {
		log.Error("Seven Ftp temp read error", err.Error())
		return
	}
	engine := database.GetMysqlEngine()
	defer engine.Close()
	for _, file := range files {
		if file.Name()[0:2] == "ol" {

			data, err := os.Open(tempDirPath + "/" + file.Name())
			if err != nil {
				log.Error("Seven Ftp tempfile read error", err.Error())
				return
			}
			scanner := bufio.NewScanner(data)
			for scanner.Scan() {
				if olSection[scanner.Text()[0:1]] == "Content" {
					runes := []byte(scanner.Text())
					var serviceType bool
					pid := string(runes[33:36])
					switch pid {
					case `851`:
						serviceType = false
					case `850`:
						serviceType = true
					default:
						serviceType = false
					}
					paymentNo := string(runes[36:44])
					payTime := string(runes[9:17]) + string(runes[23:27])
					amount := strings.TrimLeft(string(runes[76:83]), `0`)
					if amount == "" {
						amount = `0`
					}
					amountFloat, err := strconv.ParseFloat(amount, 64)
					if err != nil {
						log.Error("Seven Ol amountFloat error", err.Error())
						continue
					}
					t, _ := time.Parse(`20060102150405`, payTime+"00")
					var shipMap entity.SevenShipMapData
					var UpdateCvsShipping CvsShipping.UpdateCvsShipping
					UpdateCvsShipping.ShipType = Enum.CVS_7_ELEVEN
					flag, err := engine.Engine.Table("seven_ship_map_data").Select("paymentno_with_code,payment_no,order_id").Where("payment_no =?", paymentNo).Get(&shipMap)
					if err != nil {
						log.Error("Not find seven paymentNO "+paymentNo, err.Error())
						continue
					}
					if flag {
						var cvsShippingLogData entity.CvsShippingLogData
						flag, err := engine.Engine.Table("cvs_shipping_log_data").Select("ship_no,cvs_type,type").Where("ship_no =? && cvs_type =? && type=?", shipMap.PaymentNoWithCode, Enum.CVS_7_ELEVEN, "000").Get(&cvsShippingLogData)
						if err != nil {
							log.Error("Not find seven paymentNO "+paymentNo, err.Error())
							continue
						}
						if !flag {
							UpdateCvsShipping.ShipNo = shipMap.PaymentNoWithCode
							UpdateCvsShipping.Type = "000"
							UpdateCvsShipping.DateTime = t.Format(`2006-01-02 15:04:05`)
							UpdateCvsShipping.DetailStatus = "Y"
							UpdateCvsShipping.FlowType = "N"
							UpdateCvsShipping.Log = scanner.Text()
							UpdateCvsShipping.FileName = ""

							err = UpdateCvsShipping.UpdateCvsShippingShipment(engine)

							if err != nil {
								log.Error("Seven OL error", err)
								return
							}
							log.Info("更新訂單", shipMap.OrderId, shipMap.PaymentNoWithCode)

						}
						log.Info("seven account ol", shipMap.OrderId, shipMap.PaymentNoWithCode)
						var accountingData entity.CvsAccountingData
						var type_ = `OL`
						flag, err = engine.Engine.Table("cvs_accounting_data").Select("cvs_type,type,data_id").Where("cvs_type =? && type =? && data_id=?", Enum.CVS_7_ELEVEN, `S`, shipMap.OrderId).Get(&accountingData)
						if err != nil {
							log.Error(err.Error())
						} else {
							if !flag {
								accountingData.SetType(type_, Enum.CVS_7_ELEVEN)
								accountingData.DataId = shipMap.OrderId

								accountingData.Amount = amountFloat
								accountingData.ServiceType = serviceType
								accountingData.FileDate, _ = time.Parse(`20060102`, string(runes[9:17]))
								accountingData.Status = `1`
								accountingData.FileName = strings.Split(file.Name(), `-`)[2]
								accountingData.Log = scanner.Text()
								_ = Cvs.InsertAccountingData(engine, accountingData)

							}
						}

						log.Info("seven account ol down", shipMap.OrderId, shipMap.PaymentNoWithCode)

					}

				}
			}
		}

		log.Info("sftp connect", "7-11 OL每日出貨交寄更新結束", "Time:", time.Since(now).Seconds())
	}
	os.RemoveAll(tempDirPath)
}
func FetchCEINStatus() {
	now := time.Now()
	log.Info("sftp connect", "7-11 CEIN每日貨態更新開始", now.Format(`2006-01-02 15:04:05`))
	client := new(Client).getClient().ConnectFtp

	defer client.Close()
	files, err := client.ReadDir(folderPath["ReceiveAccept"])
	if err != nil {
		log.Error("Seven Ftp CEIN read error", err.Error())
		return
	}
	tmepDir, err := ioutil.TempDir("", "cein")
	if err != nil {
		log.Error("Seven Ftp CEIN tmep dir create error", err.Error())
		return
	}
	i := 0
	for _, file := range files {

		if i > 5 {
			break
		}
		if !file.IsDir() {

			t, _ := time.Parse(`20060102`, strings.SplitAfter(file.Name()[3:11], ".")[0])

			if t.Before(now) {
				data, err := client.Open(folderPath["ReceiveAccept"] + "/" + file.Name())

				if err != nil {
					log.Error("Seven Ftp CEIN open error"+file.Name(), err.Error())
					continue
				}
				tempfile, err := ioutil.TempFile(tmepDir, "cein-*.cein")
				if err != nil {
					log.Error("Seven Ftp CEIN open temp error"+tempfile.Name(), err.Error())
					continue
				}
				data.WriteTo(tempfile)

				err = client.Rename(folderPath["ReceiveAccept"]+"/"+file.Name(), folderPath["ReceiveAccept"]+"/done/"+file.Name())
				if err != nil {
					log.Error("Seven Ftp CEIN open temp error"+tempfile.Name(), err.Error())
					continue
				}
			}

		}
		i++
	}
	client.Close()
	files, err = ioutil.ReadDir(tmepDir)
	if err != nil {
		log.Error("Seven Ftp CEIN temp read error", err.Error())
		return
	}
	engine := database.GetMysqlEngine()
	defer engine.Close()
	for _, file := range files {

		data, err := ioutil.ReadFile(tmepDir + "/" + file.Name())
		if err != nil {
			log.Error("Seven CEIN read error", err.Error())
			continue
		}

		xmlData := CEINXml{}

		err = xml.Unmarshal(data, &xmlData)
		if err != nil {
			log.Error("Seven CEIN xml decode error", err.Error())
			continue
		}
		docDate := xmlData.DocHead.DocDate
		for _, data := range xmlData.DocContent {
			var shipMap entity.SevenShipMapData
			var UpdateCvsShipping CvsShipping.UpdateCvsShipping
			UpdateCvsShipping.ShipType = Enum.CVS_7_ELEVEN

			flag, err := engine.Engine.Table("seven_ship_map_data").Select("paymentno_with_code,payment_no").Where("payment_no =?", data.PaymentNo).Get(&shipMap)
			if err != nil {
				log.Error("Not find seven paymentNO "+data.PaymentNo, err.Error())
				continue
			}
			if flag {

				t, _ := time.Parse(`2006-01-02`, data.DCReceiveDate)
				if data.DCReceiveDate == "" {
					t, _ = time.Parse(`2006-01-02`, docDate)
				}

				UpdateCvsShipping.ShipNo = shipMap.PaymentNoWithCode

				UpdateCvsShipping.DateTime = t.Format(`2006-01-02 15:04:05`)
				UpdateCvsShipping.DetailStatus = "Y"
				if data.DCStoreStatus == "+" {
					UpdateCvsShipping.FlowType = "N"
					UpdateCvsShipping.Type = data.DCReceiveStatus

				} else if data.DCStoreStatus == "-" {
					UpdateCvsShipping.FlowType = "R"
					UpdateCvsShipping.Type = "R" + data.DCReceiveStatus

				}

				dataLog, _ := xml.Marshal(data)
				UpdateCvsShipping.Log = string(dataLog)
				UpdateCvsShipping.FileName = ""

				err = UpdateCvsShipping.OnlyWriteShippingLog(engine, false)

				if err != nil {
					log.Error("Seven CEIN error", err)
					return
				}
				log.Info("更新訂單", shipMap.OrderId, shipMap.PaymentNoWithCode)
			}
		}

	}
	log.Info("sftp connect", "7-11 CEIN每日貨態更新結束", "Time:", time.Since(now).Seconds())
}
func FetchCPPSStatus() {
	now := time.Now()

	log.Info("sftp connect", "7-11 CPPS每日貨態更新開始", now.Format(`2006-01-02 15:04:05`))

	client := new(Client).getClient().ConnectFtp

	defer client.Close()

	files, err := client.ReadDir(folderPath["ArrivedShop"])
	if err != nil {
		log.Error("Seven Ftp CPPS read error", err.Error())
		return
	}
	tmepDir, err := ioutil.TempDir("", "cpps")

	if err != nil {
		log.Error("Seven Ftp CPPS tmep dir create error", err.Error())
		return
	}
	i := 0
	for _, file := range files {

		if i > 5 {
			break
		}
		if !file.IsDir() {
			t, _ := time.Parse(`20060102`, strings.SplitAfter(file.Name()[3:11], ".")[0])

			if t.Before(now) {
				data, err := client.Open(folderPath["ArrivedShop"] + "/" + file.Name())

				if err != nil {
					log.Error("Seven Ftp CPPS open error"+file.Name(), err.Error())
					continue
				}
				tempfile, err := ioutil.TempFile(tmepDir, "cpps-*.cpps")
				if err != nil {
					log.Error("Seven Ftp CPPS open temp error"+tempfile.Name(), err.Error())
					continue
				}
				_, err = data.WriteTo(tempfile)
				if err != nil {
					log.Error("Seven Ftp CPPS write temp error"+tempfile.Name(), err.Error())
					continue
				}

				err = client.Rename(folderPath["ArrivedShop"]+"/"+file.Name(), folderPath["ArrivedShop"]+"/done/"+file.Name())
				if err != nil {
					log.Error("Seven Ftp CPPS open temp error"+tempfile.Name(), err.Error())
					continue
				}
			}

		}
		i++
	}
	client.Close()
	files, err = ioutil.ReadDir(tmepDir)
	if err != nil {
		log.Error("Seven Ftp Cpps temp read error", err.Error())
		return
	}
	engine := database.GetMysqlEngine()
	defer engine.Close()
	for _, file := range files {

		data, err := ioutil.ReadFile(tmepDir + "/" + file.Name())
		if err != nil {
			log.Error("Seven Cpps read error", err.Error())
			continue
		}

		xmlData := CPPSXml{}

		err = xml.Unmarshal(data, &xmlData)
		if err != nil {
			log.Error("Seven Cpps xml decode error", err.Error())
			continue
		}

		for _, data := range xmlData.DocContent {
			var shipMap entity.SevenShipMapData
			var UpdateCvsShipping CvsShipping.UpdateCvsShipping
			UpdateCvsShipping.ShipType = Enum.CVS_7_ELEVEN
			flag, err := engine.Engine.Table("seven_ship_map_data").Select("paymentno_with_code,payment_no").Where("payment_no =?", data.PaymentNo).Get(&shipMap)
			if err != nil {
				log.Error("Not find seven paymentNO "+data.PaymentNo, err.Error())
				continue
			}
			if flag {
				_, err := engine.Session.Table("seven_ship_map_data").Cols("ship_no").Where("payment_no =?", data.PaymentNo).Update(entity.SevenShipMapData{
					ShipNo: data.ShipmentNo,
				})
				if err != nil {
					log.Error("Not find seven paymentNO "+data.PaymentNo, err.Error())
					continue
				}
				t, _ := time.Parse(`2006-01-02150405`, data.StoreDate+data.StoreTime)

				UpdateCvsShipping.ShipNo = shipMap.PaymentNoWithCode
				UpdateCvsShipping.Type = data.StoreType
				UpdateCvsShipping.DateTime = t.Format(`2006-01-02 15:04:05`)
				UpdateCvsShipping.DetailStatus = "Y"
				UpdateCvsShipping.FlowType = "N"
				dataLog, _ := xml.Marshal(data)
				UpdateCvsShipping.Log = string(dataLog)
				UpdateCvsShipping.FileName = ""
				switch data.StoreType {
				case "101":
					err = UpdateCvsShipping.UpdateCvsShippingShop(engine)
				case "201":
					flag, err := checkCvsExist(engine, shipMap.PaymentNoWithCode, data.StoreType)
					if err != nil {
						log.Error("Not find seven paymentNO "+data.PaymentNo, err.Error())
						continue
					}
					if !flag {
						UpdateCvsShipping.FlowType = "R"
						err = UpdateCvsShipping.OnlyWriteShippingLog(engine, true)
						if err != nil {
							log.Error("Seven Cpps error", err)
							return
						}
					}

				default:
					err = UpdateCvsShipping.UpdateCvsShippingTransit(engine)
				}

				if err != nil {
					log.Error("Seven Cpps error", err)
					return
				}
				log.Info("更新訂單", shipMap.OrderId, shipMap.PaymentNoWithCode)

			}
		}

	}
	os.RemoveAll(tmepDir)

	log.Info("sftp connect", "7-11 CPPS每日貨態更新結束", "Time:", time.Since(now).Seconds())
}
func checkCvsExist(engine *database.MysqlSession, paymentNoWithCode string, storeType string) (bool, error) {
	var cvsShippingLogData entity.CvsShippingLogData
	flag, err := engine.Engine.Table("cvs_shipping_log_data").Select("ship_no,cvs_type,type").Where("ship_no =? && cvs_type =? && type=?", paymentNoWithCode, Enum.CVS_7_ELEVEN, storeType).Get(&cvsShippingLogData)
	return flag, err
}
func FetchCERTStatus() {
	now := time.Now()

	log.Info("sftp connect", "7-11 CERT每日貨態更新開始", now.Format(`20060102150405`))

	client := new(Client).getClient().ConnectFtp

	defer client.Close()
	files, err := client.ReadDir(folderPath["ReturnRequest"])
	if err != nil {
		log.Error("Seven Ftp CERT read error", err.Error())
		return
	}
	tmepDir, err := ioutil.TempDir("", "cert")

	if err != nil {
		log.Error("Seven Ftp CERT tmep dir create error", err.Error())
		return
	}
	i := 0
	for _, file := range files {
		if i > 5 {
			break
		}
		if !file.IsDir() {
			t, _ := time.Parse(`20060102`, strings.SplitAfter(file.Name()[3:11], ".")[0])

			if t.Before(now) {
				data, err := client.Open(folderPath["ReturnRequest"] + "/" + file.Name())

				if err != nil {
					log.Error("Seven Ftp  CERT open error"+file.Name(), err.Error())
					continue
				}
				tempfile, err := ioutil.TempFile(tmepDir, "cert-*.cert")
				if err != nil {
					log.Error("Seven Ftp CERT open temp error"+tempfile.Name(), err.Error())
					continue
				}
				data.WriteTo(tempfile)

				err = client.Rename(folderPath["ReturnRequest"]+"/"+file.Name(), folderPath["ReturnRequest"]+"/done/"+file.Name())
				if err != nil {
					log.Error("Seven Ftp CERT open temp error"+tempfile.Name(), err.Error())
					continue
				}
			}

		}
		i++
	}
	client.Close()
	files, err = ioutil.ReadDir(tmepDir)
	if err != nil {
		log.Error("Seven Ftp CERT temp read error", err.Error())
		return
	}
	engine := database.GetMysqlEngine()
	defer engine.Close()
	for _, file := range files {

		data, err := ioutil.ReadFile(tmepDir + "/" + file.Name())
		if err != nil {
			log.Error("Seven CERT read error", err.Error())
			continue
		}

		xmlData := CERTXml{}

		err = xml.Unmarshal(data, &xmlData)
		if err != nil {
			log.Error("Seven CERT xml decode error", err.Error())
			continue
		}

		for _, data := range xmlData.DocContent {
			var shipMap entity.SevenShipMapData
			var UpdateCvsShipping CvsShipping.UpdateCvsShipping
			UpdateCvsShipping.ShipType = Enum.CVS_7_ELEVEN
			flag, err := engine.Engine.Table("seven_ship_map_data").Select("paymentno_with_code,payment_no").Where("payment_no =?", data.PaymentNo).Get(&shipMap)
			if err != nil {
				log.Error("Not find seven paymentNO "+data.PaymentNo, err.Error())
				continue
			}
			if flag {
				var cvsShippingLogData entity.CvsShippingLogData
				t, _ := time.Parse(`2006-01-02`, data.DCPlannedReturnDate)
				returnType := "RT" + data.ReturnType
				flag, err := engine.Engine.Table("cvs_shipping_log_data").Select("ship_no,cvs_type,type").Where("ship_no =? && cvs_type =? && type=?", shipMap.PaymentNoWithCode, Enum.CVS_7_ELEVEN, returnType).Get(&cvsShippingLogData)
				if err != nil {
					log.Error("Not find seven paymentNO "+data.PaymentNo, err.Error())
					continue
				}

				if !flag {
					UpdateCvsShipping.ShipNo = shipMap.PaymentNoWithCode
					UpdateCvsShipping.Type = returnType
					UpdateCvsShipping.DateTime = t.Format(`2006-01-02 15:04:05`)
					UpdateCvsShipping.DetailStatus = "Y"
					UpdateCvsShipping.FlowType = "N"

					dataLog, _ := xml.Marshal(data)
					UpdateCvsShipping.Log = string(dataLog)
					UpdateCvsShipping.FileName = ""
					switch returnType {
					case "RT01":
						err = UpdateCvsShipping.UpdateCvsShippingBuyerNotPickUp(engine)
					default:
						err = UpdateCvsShipping.OnlyWriteShippingLog(engine, false)
					}
					if err != nil {
						log.Error("Seven CERT error", err)
						return
					}
					log.Info("更新訂單", shipMap.OrderId, shipMap.PaymentNoWithCode)
				}
			}
		}

	}
	os.RemoveAll(tmepDir)

	log.Info("sftp connect", "7-11 CERT每日貨態更新結束", "Time:", time.Since(now).Seconds())
}
func FetchCESPStatus() {
	now := time.Now()
	log.Info("sftp connect", "7-11 CESP到店取貨貨態更新開始", now.Format(`2006-01-02 15:04:05`))
	client := new(Client).getClient().ConnectFtp
	defer client.Close()
	files, err := client.ReadDir(folderPath["DailyShopOperateRecord"])
	if err != nil {
		log.Error("Seven Ftp CESP read error", err.Error())
		return
	}
	tmepDir, err := ioutil.TempDir("", "cesp")

	if err != nil {
		log.Error("Seven Ftp CESP tmep dir create error", err.Error())
		return
	}
	i := 0
	for _, file := range files {
		if i > 5 {
			break
		}
		if !file.IsDir() {
			t, _ := time.Parse(`20060102`, strings.SplitAfter(file.Name()[3:11], ".")[0])

			if t.Before(now) || t.Equal(now) {
				data, err := client.Open(folderPath["DailyShopOperateRecord"] + "/" + file.Name())

				if err != nil {
					log.Error("Seven Ftp  CESP open error"+file.Name(), err.Error())
					continue
				}
				tempfile, err := ioutil.TempFile(tmepDir, "cesp-*.cesp")
				if err != nil {
					log.Error("Seven Ftp CESP open temp error"+tempfile.Name(), err.Error())
					continue
				}
				data.WriteTo(tempfile)

				err = client.Rename(folderPath["DailyShopOperateRecord"]+"/"+file.Name(), folderPath["DailyShopOperateRecord"]+"/done/"+file.Name())
				if err != nil {
					log.Error("Seven Ftp CEDR open temp error"+tempfile.Name(), err.Error())
					continue
				}
			}

		}
		i++
	}
	client.Close()
	files, err = ioutil.ReadDir(tmepDir)
	if err != nil {
		log.Error("Seven Ftp CESP temp read error", err.Error())
		return
	}
	engine := database.GetMysqlEngine()
	defer engine.Close()

	for _, file := range files {

		data, err := ioutil.ReadFile(tmepDir + "/" + file.Name())
		if err != nil {
			log.Error("Seven CESP read error", err.Error())
			continue
		}

		xmlData := CESPXml{}

		err = xml.Unmarshal(data, &xmlData)
		if err != nil {
			log.Error("Seven CESP xml decode error", err.Error())
			continue
		}

		for _, SP := range xmlData.SP {

			for _, data := range SP.Detail {

				var shipMap entity.SevenShipMapData
				var UpdateCvsShipping CvsShipping.UpdateCvsShipping
				UpdateCvsShipping.ShipType = Enum.CVS_7_ELEVEN
				var serviceType bool
				pid := data.ParentId
				switch pid {
				case `851`:
					serviceType = false
				case `850`:
					serviceType = true
				default:
					serviceType = false
				}
				dataLog, _ := xml.Marshal(data)
				flag, err := engine.Engine.Table("seven_ship_map_data").Select("paymentno_with_code,payment_no,order_id").Where("payment_no =?", data.PaymentNo).Get(&shipMap)
				if err != nil {
					log.Error("Not find seven paymentNO "+data.PaymentNo, err.Error())
					continue
				}
				if flag {

					var cvsShippingLogData entity.CvsShippingLogData
					t, _ := time.Parse(`2006-01-02`, data.SPDate)
					flowType := "777"
					if data.DCStoreStatus == "-" {
						flowType = "888"
					}
					flag, err := engine.Engine.Table("cvs_shipping_log_data").Select("ship_no,cvs_type,type").Where("ship_no =? && cvs_type =? && type=?", shipMap.PaymentNoWithCode, Enum.CVS_7_ELEVEN, flowType).Get(&cvsShippingLogData)
					if err != nil {
						log.Error("Not find seven paymentNO "+data.PaymentNo, err.Error())
						continue
					}
					if !flag {

						UpdateCvsShipping.ShipNo = shipMap.PaymentNoWithCode
						UpdateCvsShipping.Type = flowType
						UpdateCvsShipping.DateTime = t.Format(`2006-01-02 15:04:05`)
						UpdateCvsShipping.DetailStatus = "Y"
						UpdateCvsShipping.FlowType = "N"

						UpdateCvsShipping.Log = string(dataLog)
						UpdateCvsShipping.FileName = ""
						switch flowType {
						case "777":
							err = UpdateCvsShipping.UpdateCvsShippingSuccess(engine)
						case "888":
							err = UpdateCvsShipping.OnlyWriteShippingLog(engine, true)
						}

						if err != nil {
							log.Error("Seven CESP error", err)
							return
						}
						log.Info("更新訂單", shipMap.OrderId, shipMap.PaymentNoWithCode)

					}
					log.Info("seven account send", shipMap.OrderId, shipMap.PaymentNoWithCode)
					var accountingData entity.CvsAccountingData
					var type_ = `CESP`
					flag, err = engine.Engine.Table("cvs_accounting_data").Select("cvs_type,type,data_id").Where("cvs_type =? && type =? && data_id=?", Enum.CVS_7_ELEVEN, `P`, shipMap.OrderId).Get(&accountingData)
					if err != nil {
						log.Error(err.Error())
					} else {
						if !flag {

							accountingData.SetType(type_, Enum.CVS_7_ELEVEN)
							accountingData.DataId = shipMap.OrderId
							accountingData.Amount, _ = strconv.ParseFloat(data.SPAmount, 64)
							accountingData.ServiceType = serviceType
							accountingData.FileDate, _ = time.Parse(`2006-01-02`, data.SPDate)
							accountingData.Status = `1`
							accountingData.FileName = file.Name()
							accountingData.Log = string(dataLog)
							_ = Cvs.InsertAccountingData(engine, accountingData)

						}
					}
					log.Info("seven account send end", shipMap.OrderId, shipMap.PaymentNoWithCode)
				}
			}
		}
	}
	os.RemoveAll(tmepDir)

	log.Info("sftp connect", "7-11 CESP到店取貨貨態更新結束", "Time:", time.Since(now).Seconds())
}
func FetchCEDRStatus() {
	now := time.Now()

	log.Info("sftp connect", "7-11 CEDR退貨貨態更新開始", now.Format(`2006-01-02 15:04:05`))
	client := new(Client).getClient().ConnectFtp

	defer client.Close()
	files, err := client.ReadDir(folderPath["ReturnAccept"])
	if err != nil {
		log.Error("Seven Ftp CEDR read error", err.Error())
		return
	}
	tmepDir, err := ioutil.TempDir("", "cedr")

	if err != nil {
		log.Error("Seven Ftp CEDR tmep dir create error", err.Error())
		return
	}
	i := 0
	for _, file := range files {
		if i > 5 {
			break
		}
		if !file.IsDir() {
			t, _ := time.Parse(`20060102`, strings.SplitAfter(file.Name()[3:11], ".")[0])

			if t.Before(now) {
				data, err := client.Open(folderPath["ReturnAccept"] + "/" + file.Name())

				if err != nil {
					log.Error("Seven Ftp  CEDR open error"+file.Name(), err.Error())
					continue
				}
				tempfile, err := ioutil.TempFile(tmepDir, "cedt-*.cedt")
				if err != nil {
					log.Error("Seven Ftp CEDR open temp error"+tempfile.Name(), err.Error())
					continue
				}
				data.WriteTo(tempfile)

				err = client.Rename(folderPath["ReturnAccept"]+"/"+file.Name(), folderPath["ReturnAccept"]+"/done/"+file.Name())
				if err != nil {
					log.Error("Seven Ftp CEDR open temp error"+tempfile.Name(), err.Error())
					continue
				}
			}

		}
		i++
	}
	client.Close()
	files, err = ioutil.ReadDir(tmepDir)
	if err != nil {
		log.Error("Seven Ftp CEDR temp read error", err.Error())
		return
	}
	engine := database.GetMysqlEngine()
	defer engine.Close()
	for _, file := range files {

		data, err := ioutil.ReadFile(tmepDir + "/" + file.Name())
		if err != nil {
			log.Error("Seven CEDR read error", err.Error())
			continue
		}

		xmlData := CEDRXml{}

		err = xml.Unmarshal(data, &xmlData)
		if err != nil {
			log.Error("Seven CEDR xml decode error", err.Error())
			continue
		}

		for _, data := range xmlData.DocContent {
			var shipMap entity.SevenShipMapData
			var UpdateCvsShipping CvsShipping.UpdateCvsShipping
			UpdateCvsShipping.ShipType = Enum.CVS_7_ELEVEN
			flag, err := engine.Engine.Table("seven_ship_map_data").Select("paymentno_with_code,payment_no").Where("payment_no =?", data.PaymentNo).Get(&shipMap)
			if err != nil {
				log.Error("Not find seven paymentNO "+data.PaymentNo, err.Error())
				continue
			}
			if flag {

				var cvsShippingLogData entity.CvsShippingLogData
				t, _ := time.Parse(`2006-01-02`, data.DCReturnDate)
				returnCode := "DR" + data.DCReturnCode
				flag, err := engine.Engine.Table("cvs_shipping_log_data").Select("ship_no,cvs_type,type").Where("ship_no =? && cvs_type =? && type=? && is_show=?", shipMap.PaymentNoWithCode, Enum.CVS_7_ELEVEN, returnCode, 0).Get(&cvsShippingLogData)
				if err != nil {
					log.Error("Not find seven paymentNO "+data.PaymentNo, err.Error())
					continue
				}
				if !flag {

					UpdateCvsShipping.ShipNo = shipMap.PaymentNoWithCode
					UpdateCvsShipping.Type = returnCode
					UpdateCvsShipping.DateTime = t.Format(`2006-01-02 15:04:05`)
					UpdateCvsShipping.DetailStatus = "Y"
					UpdateCvsShipping.FlowType = "N"
					dataLog, _ := xml.Marshal(data)
					UpdateCvsShipping.Log = string(dataLog)
					UpdateCvsShipping.FileName = ""

					err = UpdateCvsShipping.OnlyWriteShippingLog(engine, false)

					if err != nil {
						log.Error("Seven CEDR error", err)
						return
					}
					log.Info("更新訂單", shipMap.OrderId, shipMap.PaymentNoWithCode)
				}
			}
		}

	}
	os.RemoveAll(tmepDir)

	log.Info("sftp connect", "7-11 CEDR退貨貨態更新結束", "Time", time.Since(now).Seconds())
}

//FetchDailyShopStatus is update seven shops by daily ftp file
func FetchDailyShopStatus() {
	var enc = traditionalchinese.Big5
	nowTime := time.Now().Format("20060102")
	log.Info("sftp connect", "7-11每日更新店鋪開始", nowTime)
	fileName := "01" + nowTime
	client := new(Client).getClient()

	defer client.ConnectFtp.Close()

	srcfile, err := client.ConnectFtp.Open(folderPath["Shop"] + "/" + fileName + ".STD")

	if err != nil {
		log.Error(nowTime + "seven 店鋪取得失敗 " + err.Error())
		return
	}
	tempfile, err := ioutil.TempFile("", "ctsd-*.STD")
	if err != nil {
		log.Error(nowTime + "seven 店鋪取得失敗 " + err.Error())
		return
	}

	srcfile.WriteTo(tempfile)
	client.ConnectFtp.Close()
	datafile, err := os.Open(tempfile.Name())
	// Read UTF-8 from a GBK encoded file.
	if err != nil {
		log.Error(nowTime + "seven 店鋪取得失敗 " + err.Error())
		return
	}
	r := transform.NewReader(datafile, enc.NewDecoder())

	scanner := bufio.NewScanner(r)
	n := 0
	shops := []entity.SevenMyshipShopData{}
	for scanner.Scan() {

		if len(scanner.Text()) > 100 {
			var data entity.SevenMyshipShopData

			if n == 0 {
				s := strings.Replace(scanner.Text(), "\uFEFF", "", -1)
				data, flag = DataSplit(standardizeSpaces(s))
			} else {
				data, flag = DataSplit(standardizeSpaces(scanner.Text()))
			}
			if flag {
				shops = append(shops, data)
			}
		}

		n++
	}
	tempfile.Close()
	os.Remove(tempfile.Name())
	sevenmyshipdao.InsertAddressByDailyFile(shops)
	sevenmyshipdao.UpdateClosedShop()
	log.Info("sftp connect", "7-11每日更新店鋪結束", nowTime)
}
func DataSplit(address string) (entity.SevenMyshipShopData, bool) {
	k := strings.Split(standardizeSpaces(address), " ")
	var flag bool
	data := new(entity.SevenMyshipShopData)

	data.StoreName = string([]rune(k[0])[6:len([]rune(k[0]))])
	data.StoreID = SubStr(k[0], 0, 5)
	data.Country, data.District, data.Address, flag = AddressSplit(k[2][6:])
	data.Opened = true
	return *data, flag
}

func AddressSplit(address string) (string, string, string, bool) {

	if country, ok := FindCountry(address); ok {

		subAddress := strings.TrimPrefix(address, country)

		if district, ok := FindDistrictByCountry(subAddress, country); ok {

			return country, district, strings.TrimPrefix(subAddress, district), true
		}
		return "", "", "", false

	}
	return "", "", "", false

}
func FindCountry(address string) (string, bool) {

	for _, country := range Enum.CountryList {

		clen := utf8.RuneCountInString(country)
		subAddress := SubStr(address, 0, clen-1)
		if subAddress == country {
			return country, true
		}
	}
	return "", false
}
func FindDistrictByCountry(address string, country string) (string, bool) {

	if _, ok := Enum.CountryCityList[country]; ok {
		for _, district := range Enum.CountryCityList[country] {
			clen := utf8.RuneCountInString(district)
			subAddress := SubStr(address, 0, clen-1)
			if subAddress == district {
				return district, true
			}
		}
		return "", false
	}
	return "", false
}
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func SubStr(source interface{}, start int, end int) string {
	str := source.(string)
	var r = []rune(str)
	length := len(r)
	subLen := end - start

	for {
		if start < 0 {
			break
		}
		if start == 0 && subLen == length {
			break
		}
		if end > length {
			subLen = length - start
		}
		if end < 0 {
			subLen = length - start + end
		}
		var substring bytes.Buffer
		if end > 0 {
			subLen = subLen + 1
		}
		for i := start; i < subLen; i++ {
			substring.WriteString(string(r[i]))
		}
		str = substring.String()

		break
	}

	return str
}
