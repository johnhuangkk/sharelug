package PostBag

import (
	"api/services/Enum"
	"api/services/Service/Balance"
	"api/services/Service/OrderService"
	"api/services/VO/ShipmentVO"
	"api/services/dao/Orders"
	postBag "api/services/dao/PostBag"
	"api/services/dao/iPost"
	"api/services/dao/sequence"
	"api/services/database"
	"api/services/entity"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/tools"
	sxml "api/services/util/xml"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/viper"
)

type PostBagXml struct {
	XMLName xml.Name    `xml:"LgsMsg"`
	Header  PostBagHead `xml:"Head"`
	Datas   []PostBagData
}
type PostBagHead struct {
	VipNo     string
	MailType  string
	ValidDate string
	DataCnt   string
}
type PostBagData struct {
	XMLName    xml.Name `xml:"Data"`
	MailNo     string
	GiroAcntNo string
	CodAmt     string
}
type address struct {
	Name, Mobile, Zipcode, Alias, Address string
}
type PostBagConsignmentData struct {
	MerchantId, MerchantName, ShopName, ShipNumber, Mark, Type string
	Seller, Receiver                                           address
}

func CreateShipNumber(engine *database.MysqlSession, order entity.OrderData, seller entity.MemberData, sellerAddress *ShipmentVO.SellerSenderAddress) (string, error) {
	seq, err := sequence.GetPostBagSeq()

	if err != nil {
		log.Error(order.OrderId, "PostBag Seq Create Fail")
		return order.OrderId, fmt.Errorf("%s", "系統忙碌中，請稍後再進行出貨。")
	}
	packageOfficeNumber := viper.GetString(`POSTBAG.packageNumber`)
	receiverPostCode := strings.Split(order.ReceiverAddress, ",")
	postCode := receiverPostCode[0] + "00"
	shipNumber := tools.StringPadLeft(seq, 6) + packageOfficeNumber + "18" + postCode
	if len(shipNumber) != 19 {
		log.Error(order.OrderId, shipNumber, "PostBag Seq Create Fail , Length too short", len(shipNumber))
		return order.OrderId, fmt.Errorf("%s", "系統忙碌中，請稍後再進行出貨。")
	}
	check, err := calculateCheckSum(shipNumber)
	if err != nil {
		log.Error(order.OrderId, shipNumber, "PostBag Seq Create Fail, CheckSumFail")
		return order.OrderId, fmt.Errorf("%s", "系統忙碌中，請稍後再進行出貨。")
	}
	shipNumber += check
	var postbagInfo = entity.PostBagConsignmentData{}
	postbagInfo.OrderId = order.OrderId
	postbagInfo.MerchantId = viper.GetString(`POSTBAG.officialNumber`)
	postbagInfo.ShipNumber = shipNumber
	postbagInfo.SellerId = seller.Uid
	postbagInfo.SellerName = seller.SendName
	postbagInfo.SellerPhone = seller.Mphone
	postbagInfo.SellerZip = sellerAddress.Zip + "00"
	postbagInfo.SellerAddr = sellerAddress.Address
	postbagInfo.VerifyFileName = ""

	_, err = postBag.InsertPostBagConsignmentData(engine, postbagInfo)
	if err != nil {
		log.Error(err.Error(), "PostBag Data Create Fail")
		return "", err
	}
	return shipNumber, nil
}

func calculateCheckSum(shipNumber string) (string, error) {
	odd := 0
	even := 0
	numArr := strings.Split(shipNumber, "")
	for i, v := range numArr {
		num, _ := strconv.Atoi(v)
		if i%2 == 0 {
			even += num
		} else {
			odd += num
		}
	}
	sum := 10 - (odd*3+even)%10
	if sum == 10 {
		return strconv.Itoa(0), nil
	}

	return strconv.Itoa(sum), nil
}

func GetConsignment(orderIds []string) ([]PostBagConsignmentData, error) {
	var data []PostBagConsignmentData
	consignments, err := postBag.GetConsignmentData(orderIds)
	if err != nil {
		log.Error("Get PostBag Consignment Error", err.Error())
		return data, err
	}

	for _, consignment := range consignments {
		receiverAddrInfo := strings.Split(consignment.Order.ReceiverAddress, ",")
		bagData := PostBagConsignmentData{}
		bagData.ShipNumber = consignment.Order.ShipNumber
		bagData.MerchantId = consignment.PostBag.MerchantId
		bagData.MerchantName = "Sharelug"
		bagData.ShopName = consignment.Store.StoreName
		bagData.Seller.Address = consignment.PostBag.SellerAddr
		bagData.Seller.Zipcode = consignment.PostBag.SellerZip
		bagData.Seller.Name = consignment.PostBag.SellerName
		bagData.Seller.Mobile = consignment.PostBag.SellerPhone
		bagData.Receiver.Address = strings.Join(receiverAddrInfo[1:], "")
		bagData.Receiver.Zipcode = receiverAddrInfo[0] + "00"
		bagData.Receiver.Mobile = consignment.Order.ReceiverPhone
		bagData.Receiver.Name = consignment.Order.ReceiverName
		switch consignment.Order.ShipType {
		case Enum.DELIVERY_POST_BAG1:
			bagData.Type = "1"
		case Enum.DELIVERY_POST_BAG2:
			bagData.Type = "2"
		case Enum.DELIVERY_POST_BAG3:
			bagData.Type = "3"
		}
		data = append(data, bagData)
	}
	return data, nil
}
func CheckFileUpdate() {
	nTime := time.Now()
	log.Info("PostBag Shipping Number CheckUpdate Start", time.Now().Format(`2006-01-02 15:04:05`))
	// consignments, err := postBag.FindRecentNonVerifyFilePostBagConsignmentData()
	// //連線拿清單比對
	c, err := getPostClient()
	if err != nil {
		log.Error("PostBag FTP connect fail", err)
		return
	}
	defer c.Quit()
	c.ChangeDir("/Result")
	s, err := c.NameList("/Result")
	if err != nil {
		log.Error(err.Error(), "PostBag Shipping Number CheckUpdate Fail")
		return
	}
	c.Quit()
	rfiles, err := ioutil.ReadDir("./data/postbag/done/")

	engine := database.GetMysqlEngine()
	defer engine.Close()
	if len(rfiles) > 0 {

		err = engine.Session.Begin()
		if err != nil {
			log.Error(err.Error(), "PostBag Shipping Number CheckUpdate Fail")
			return
		}
		for _, rfile := range rfiles {
			if !rfile.IsDir() {
				fileName := strings.Split(rfile.Name(), ".")[0]
				for _, name := range s {
					rName := strings.Split(name, "_")[0]
					if rName == fileName {
						_, err := engine.Session.Table(entity.PostBagConsignmentData{}).Where(`file_name =?`, fileName+".xml").Update(&entity.PostBagConsignmentData{
							VerifyFileName: name,
						})
						if err != nil {
							log.Error(err.Error(), "PostBag Shipping Number CheckUpdate update data Fail")
							return
						}
					}
				}
			}
		}
		err = engine.Session.Commit()
		if err != nil {
			log.Error(err.Error(), "PostBag Shipping Number CheckUpdate Fail")
			return
		}
	}

	var datas []entity.PostBagConsignmentData
	err = engine.Engine.Table(entity.PostBagConsignmentData{}).
		Where(`post_bag_consignment_data.file_name != ?`).
		Where(`post_bag_consignment_data.verify_file_name =?`, "").
		Find(&datas)
	if err != nil {
		log.Error(err.Error(), "PostBag Shipping Number CheckUpdate Fail")
		return
	}
	if len(datas) > 0 {
		for _, data := range datas {
			if int(nTime.Sub(data.VerifyTime).Minutes()) >= 25 {
				_, err = engine.Session.Table(entity.PostBagConsignmentData{}).Where(`order_id=?`, data.OrderId).Cols("file_name").Update(entity.PostBagConsignmentData{
					FileName: "",
				})
				if err != nil {
					log.Error(err.Error(), "PostBag Shipping Number CheckUpdate Fail")
					continue
				}
			}
		}
	}
	if err != nil {
		log.Error(err.Error(), "PostBag Shipping Number CheckUpdate update data Fail")
		return
	}
	log.Info("PostBag Shipping Number CheckUpdate Finish", time.Now().Format(`2006-01-02 15:04:05`))
}
func BuildXml() {
	log.Info("PostBag Shipping Number Build Xml Start", time.Now().Format(`2006-01-02 15:04:05`))
	consignments, err := postBag.FindRecentNonVerifyPostBagConsignmentData()

	if err != nil {
		log.Error(err.Error(), "PostBag Shipping Number Build Xml Fail")
		return
	}
	if len(consignments) > 0 {
		nowT := time.Now()
		postBagXml := &PostBagXml{}
		postBagXml.Header.VipNo = viper.GetString(`POSTBAG.officialNumber`)
		postBagXml.Header.MailType = "18"
		postBagXml.Header.DataCnt = strconv.Itoa(len(consignments))
		postBagXml.Header.ValidDate = nowT.AddDate(0, 0, 7).Format("20060102")
		var numbers []string
		for _, consignment := range consignments {
			data := PostBagData{}
			data.MailNo = consignment.PostBag.ShipNumber
			data.GiroAcntNo = "83513049"
			data.CodAmt = "0"
			postBagXml.Datas = append(postBagXml.Datas, data)
			numbers = append(numbers, consignment.PostBag.ShipNumber)
		}

		filename := "Val" + viper.GetString(`POSTBAG.officialNumber`) + nowT.Format("20060102150405") + ".xml"
		err = postBag.UpdatePostBagConsignmentData(numbers, filename, nowT)
		if err != nil {
			log.Error(err.Error(), "update error")
			return
		}

		out, err := os.OpenFile("./data/postbag/"+filename, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Error(err.Error(), "PostBag Shipping Number Build Xml Fail ,os error")
		}
		out.WriteString(xml.Header)

		data, err := sxml.PostBagXmlEncoder(&postBagXml)
		if err != nil {
			log.Error(err.Error(), "PostBag Shipping Number Build Xml Fail")
			return
		}
		path := tools.GetFilePath("/postbag/", "", 0)
		out.Close()
		_, err = tools.CreateFile(path, string(data), filename)
		if err != nil {
			log.Error(err.Error(), "PostBag Shipping Number Build Xml Fail")
			return
		}

	}

	text := fmt.Sprintf("PostBag Shipping Number Build Xml Finish,Total %d shippment", len(consignments))
	log.Info(text, time.Now().Format(`2006-01-02 15:04:05`))

}
func UploadShipOrderFile() {
	log.Info("PostBag Shipping Number Xml Upload Start", time.Now().Format(`2006-01-02 15:04:05`))
	rfiles, err := ioutil.ReadDir("./data/postbag/")

	if err != nil {
		log.Error("PostBag xml data folder read fail", err.Error())
		return
	}
	var uploadedFiles []string
	if len(rfiles) > 1 {

		c, err := getPostClient()
		if err != nil {
			log.Error("PostBag FTP connect fail", err)
			return
		}

		defer c.Quit()

		for _, rfile := range rfiles {
			if !rfile.IsDir() {
				path := tools.GetFilePath("/postbag/", "", 0)
				file, err := os.Open(path + rfile.Name())

				if err != nil {
					log.Error(err.Error(), "PostBag FTP File Error")
					continue
				}
				defer file.Close()

				fileNames := strings.Split(file.Name(), "/")
				err = c.Stor(fileNames[3], file)
				if err != nil {
					log.Error("PostBag FTP Upload Fail", err)
					return
				}
				uploadedFiles = append(uploadedFiles, file.Name())
			}
		}
		for _, name := range uploadedFiles {

			names := strings.Split(name, "/")

			err := os.Rename(name, "./data/postbag/done/"+names[3])
			if err != nil {
				log.Error(err.Error())
				continue
			}
		}
	}
	text := fmt.Sprintf("PostBag Shipping Number Xml Upload Finish,Total %d Files", len(uploadedFiles))
	log.Info(text, time.Now().Format(`2006-01-02 15:04:05`))
}
func getPostClient() (*ftp.ServerConn, error) {
	c, err := ftp.Dial(viper.GetString(`POSTBAG.host`), ftp.DialWithTimeout(5*time.Second))
	if err != nil {

		return c, err
	}

	err = c.Login(viper.GetString(`POSTBAG.id`), viper.GetString(`POSTBAG.pwd`))
	if err != nil {

		return c, err
	}
	return c, nil
}
func getClient() (*ftp.ServerConn, error) {
	c, err := ftp.Dial(viper.GetString(`POSTBAG.sharelugHost`), ftp.DialWithTimeout(5*time.Second))
	if err != nil {

		return c, err
	}
	err = c.Login(viper.GetString(`POSTBAG.sharelugid`), viper.GetString(`POSTBAG.sharelugpw`))
	if err != nil {

		return c, err
	}
	return c, nil
}
func UpdateStatus() {

	log.Info("PostBag Shipping Status Update", time.Now().Format(`2006-01-02 15:04:05`))
	var err error
	var c *ftp.ServerConn

	c, err = getClient()
	if err != nil {
		log.Error("PostBag FTP connect fail", err.Error())
		return
	}
	defer c.Quit()
	vs, err := c.NameList(`status`)
	if err != nil {
		log.Error("PostBag FTP shipping status read fail", err.Error())
		return
	}
	for _, v := range vs {

		r, err := c.Retr(v)

		if err != nil {
			log.Error("PostBag FTP shipping status read fail", err.Error())
			continue
		}

		defer r.Close()

		buf, err := ioutil.ReadAll(r)

		if err != nil {
			log.Error("PostBag FTP shipping status read fail", err.Error())
			continue
		}
		r.Close()
		lines := strings.Split(string(buf), "\n")
		for i, line := range lines {
			if i == 0 {
				line = strings.Replace(line, "\uFEFF", "", -1)
			}
			contexts := strings.Split(line, "|")
			if len(contexts) >= 5 {

				InsertInfo(contexts)

			}

		}
		c.Rename(v, "done/"+v)
	}
	log.Info("PostBag Shipping Status Finish", time.Now().Format(`2006-01-02 15:04:05`))
}

func InsertInfo(lineInfo []string) {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var postShippingStatus = &entity.PostShippingStatus{}
	var orderData entity.OrderData

	t, _ := time.Parse(`20060102150405`, lineInfo[3])

	postShippingStatus.MailNo = lineInfo[0]
	postShippingStatus.HandleTime = t.Format(`2006-01-02 15:04:05`)
	postShippingStatus.ShippingStatus = lineInfo[5]
	postShippingStatus.Branch = lineInfo[2]
	postShippingStatus.StatusCode = lineInfo[4]
	postShippingStatus.CreateTime = tools.Now(`YmdHis`)
	postShippingStatus.Detail = strings.Join(lineInfo, ",")

	orderData, _ = Orders.GetOrderDataByPostBag(engine, postShippingStatus.MailNo)

	switch postShippingStatus.StatusCode {
	case `A200`:
		if len(orderData.OrderId) > 0 {
			if orderData.ShipStatus == Enum.OrderShipTake {
				OrderService.OrderCaptureRelease(&orderData, time.Time{})
				_ = Balance.OrderShipDeduction(engine, &orderData)
				// 回寫訂單出貨時間
				orderData.ShipTime = t
				_ = model.UpdateOrderDataShipStatus(engine, orderData, Enum.OrderShipment)
			}
		}
	case `Z400`:
		if len(orderData.OrderId) > 0 {
			if orderData.ShipStatus != Enum.OrderShipTransit && orderData.ShipStatus != Enum.OrderSuccess {
				_ = model.UpdateOrderDataShipStatus(engine, orderData, Enum.OrderShipTransit)
			}
		}
	// case `Y400`:
	// 	if len(orderData.OrderId) > 0 {
	// 		if orderData.ShipStatus != Enum.OrderShipTransit && orderData.ShipStatus != Enum.OrderSuccess {
	// 			_ = model.UpdateOrderDataShipStatus(engine, orderData, Enum.OrderShipTransit)
	// 		}
	// 	}
	case `I400`:
		if len(orderData.OrderId) > 0 {
			if orderData.ShipStatus != Enum.OrderSuccess {
				_ = model.UpdateOrderDataShipStatus(engine, orderData, Enum.OrderSuccess)
			}
		}
	case `I500`:
		if len(orderData.OrderId) > 0 {
			if orderData.ShipStatus != Enum.OrderSuccess {
				_ = model.UpdateOrderDataShipStatus(engine, orderData, Enum.OrderSuccess)
			}
		}
	}
	// if postShippingStatus.ShippingStatus == `郵局貨件轉運中` || postShippingStatus.StatusCode == `Z400` {

	// }

	// else if postShippingStatus.ShippingStatus == `貨件投遞中` || postShippingStatus.StatusCode == `Y400` {
	// 	if len(orderData.OrderId) > 0 {
	// 		if orderData.ShipStatus != Enum.OrderShipTransit && orderData.ShipStatus != Enum.OrderSuccess {
	// 			_ = model.UpdateOrderDataShipStatus(engine, orderData, Enum.OrderShipTransit)
	// 		}
	// 	}
	// }

	// else if postShippingStatus.StatusCode == `I400` || postShippingStatus.ShippingStatus == `投遞成功` {
	// 	if len(orderData.OrderId) > 0 {
	// 		if orderData.ShipStatus != Enum.OrderSuccess {
	// 			_ = model.UpdateOrderDataShipStatus(engine, orderData, Enum.OrderSuccess)
	// 		}
	// 	}
	// }
	err := iPost.InsertPostBagShippingStatus(engine, *postShippingStatus)
	if err != nil {
		log.Error("InsertPostBagShippingStatus fail %v", postShippingStatus)
	}
}
