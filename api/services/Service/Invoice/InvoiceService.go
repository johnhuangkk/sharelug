package Invoice

import (
	"api/services/Enum"
	"api/services/Service/Invoice/InvoiceXml"
	"api/services/Service/Mail"
	"api/services/Service/MemberService"
	"api/services/Service/Notification"
	"api/services/VO/InvoiceVo"
	"api/services/VO/Request"
	"api/services/dao/InvoiceDao"
	"api/services/dao/Orders"
	"api/services/dao/member"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/tools"
	"api/services/util/validate"
	"api/services/util/xml"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func CancelAllowanceExample(allowanceId, reason string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	data, err := InvoiceDao.GetAllowanceDataByAllowanceId(engine, allowanceId)
	if err != nil {
		log.Error("Get Allowance Data Error", err)
		return err
	}
	data.AllowanceStatus = Enum.AllowanceStatusCancel
	data.CancelTime = time.Now()
	data.CancelReason = reason
	if err := InvoiceDao.UpdateAllowanceData(engine, data); err != nil {
		log.Error("Update Allowance Data Error", err)
		return err
	}
	if err := cancelAllowance(data); err != nil {
		log.Error("create Allowance Xml Error", err)
		return err
	}
	return nil
}

func AllowanceExample(orderId string) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	order, _ := InvoiceDao.GetInvoiceByOrderId(engine, orderId)
	var vo InvoiceVo.AllowanceVo
	vo.Buyer = order.Buyer
	vo.Identifier = order.Identifier
	var items []InvoiceVo.Details
	log.Debug("sss", order)
	if err := json.Unmarshal([]byte(order.Detail), &items); err != nil {
		log.Error("Unmarshal Error", err)
		return err
	}
	var details []InvoiceVo.ProductItemVo
	for k, v := range items {
		detail := InvoiceVo.ProductItemVo {
			OriginalInvoiceNumber: fmt.Sprintf("%s%s", order.InvoiceTrack, order.InvoiceNumber),
			OriginalDescription: v.ProductName,
			OriginalInvoiceDate: order.CreateTime.Format("20060102"),
			OriginalSequenceNumber: strconv.Itoa(int(v.Sequence)),
			Quantity: v.Quantity,
			UnitPrice: v.ProductPrice,
			Amount: v.ProductAmount,
			Tax: int64(tools.Round(float64(v.ProductAmount) * Enum.InvoiceTaxRate)),
			TaxType: Enum.InvoiceTaxType01,
			AllowanceSequenceNumber: strconv.Itoa(k+1),
		}
		details = append(details, detail)
	}
	vo.Products = details
	if err := ProcessCreateAllowance(engine, vo);err != nil {
		log.Error("Create Allowance Error", err)
		return err
	}
	return nil
}
//??????C2C??????(?????????????????????)
func ProcessCreateServiceInvoice(engine *database.MysqlSession, OrderId string) error {
	result, _ := Orders.GetOrderByOrderId(engine, OrderId)
	var items []InvoiceVo.Details
	fee := result.PlatformShipFee + result.PlatformTransFee + result.PlatformInfoFee + result.PlatformPayFee
	item := InvoiceVo.Details {
		ProductName: "???????????????",
		Quantity: 1,
		ProductPrice: int64(fee),
		ProductAmount: int64(fee),
		Sequence: int64(1),
	}
	items = append(items, item)
	user, err := member.GetMemberDataByUid(engine, result.SellerId)
	if err != nil {
		log.Error("Get User Data Error", err)
		return err
	}
	var carrier Request.CarrierRequest
	//??????????????????????????????
	if user.Category != Enum.CategoryCompany {
		data, err := MemberService.GetMemberCarrierByMemberId(engine, result.SellerId)
		if err != nil {
			log.Error("Get Member Carrier Data Error", err)
			return err
		}
		carrier.InvoiceType = data.InvoiceType
		carrier.CompanyName = data.CompanyName
		carrier.CompanyBan = data.CompanyBan
		carrier.DonateBan = data.DonateBan
		carrier.CarrierType = data.CarrierType
		carrier.CarrierId = data.CarrierId
	} else {
		Key := viper.GetString("EncryptKey")
		carrier.InvoiceType = Enum.InvoiceTypeCompany
		carrier.CompanyName = user.CompanyName
		carrier.CompanyBan = tools.AesDecrypt(user.Identity, Key)
		carrier.DonateBan = ""
		carrier.CarrierType = Enum.InvoiceCarrierTypeMember
		carrier.CarrierId = user.InvoiceCarrier
	}
	vo := InvoiceVo.InvoiceVo{
		OrderId: result.OrderId,
		Details: items,
		UserId: result.SellerId,
		Amount: int64(fee),
		Carrier: carrier,
	}
	data, err := CreateInvoice(engine, vo)
	if err != nil {
		log.Error("Create Invoice Data Error", err)
		return err
	}
	if err := Notification.SendOrderInvoiceMessage(engine, user, data);err != nil {
		return err
	}
	return nil
}
//??????B2C??????(??????????????????)
func ProcessCreateB2cInvoice(engine *database.MysqlSession, OrderId string) error {
	result, _ := Orders.GetB2cOrderByOrderId(engine, OrderId)
	var details []entity.B2cOrderDetail
	if err := tools.JsonDecode([]byte(result.OrderDetail), &details); err != nil {
		log.Error("Unmarshal Error", err)
		return err
	}
	var items []InvoiceVo.Details
	for k, v := range details {
		item := InvoiceVo.Details {
			ProductName: v.ProductName,
			Quantity: 1,
			ProductPrice: v.ProductAmount,
			ProductAmount: v.ProductAmount,
			Sequence: int64(k+1),
		}
		items = append(items, item)
	}
	vo := InvoiceVo.InvoiceVo{
		OrderId: result.OrderId,
		Details: items,
		UserId: result.UserId,
		Amount: result.Amount,
		Carrier: Request.CarrierRequest {
			InvoiceType: result.InvoiceType,
			CompanyName: result.CompanyName,
			CompanyBan: result.CompanyBan,
			DonateBan: result.DonateBan,
			CarrierType: result.CarrierType,
			CarrierId: result.CarrierId,
		},
	}
	user, err := member.GetMemberDataByUid(engine, result.UserId)
	if err != nil {
		log.Error("Get User Data Error", err)
		return err
	}
	data, err := CreateInvoice(engine, vo)
	if err != nil {
		log.Error("Create Invoice Data Error", err)
		return err
	}
	if err := Notification.SendServiceInvoiceMessage(engine, user, data);err != nil {
		return err
	}
	return nil
}
//????????????
func ProcessCancelInvoice(engine *database.MysqlSession, OrderId, reason string) error {
	data, err := InvoiceDao.GetInvoiceByOrderId(engine, OrderId)
	if err != nil {
		log.Error("Get Invoice Data Error", err)
		return err
	}
	data.InvoiceStatus = Enum.InvoiceStatusCancel
	data.CancelReason = reason
	data.CancelTime = time.Now()
	if err := InvoiceDao.UpdateInvoiceData(engine, data); err != nil {
		log.Error("Update Invoice Data Error", err)
		return err
	}
	if err := CancelInvoice(data); err != nil {
		log.Error("Cancel Invoice Xml Error", err)
		return err
	}
	return nil
}
//????????????
func ProcessVoidInvoice(engine *database.MysqlSession, OrderId, reason string) error {
	data, err := InvoiceDao.GetInvoiceByOrderId(engine, OrderId)
	if err != nil {
		log.Error("Get Invoice Data Error", err)
		return err
	}
	data.InvoiceStatus = Enum.InvoiceStatusCancel
	data.VoidReason = reason
	data.VoidTime = time.Now()
	if err := InvoiceDao.UpdateInvoiceData(engine, data); err != nil {
		log.Error("Update Invoice Data Error", err)
		return err
	}
	if err := VoidInvoice(data); err != nil {
		log.Error("Void Invoice Xml Error", err)
		return err
	}
	return nil
}
//???????????????
func ProcessCreateAllowance(engine *database.MysqlSession, vo InvoiceVo.AllowanceVo) error {
	taxAmount := int64(0)
	totalAmount := int64(0)
	for _, v := range vo.Products {
		taxAmount += v.Tax
		totalAmount += v.Amount
	}
	jsonStr, _ := json.Marshal(vo.Products)
	data := entity.AllowanceData{
		AllowanceId: tools.GeneratorAllowanceId(), //????????????
		AllowanceDate: time.Now().Format("20060102"),
		AllowanceType: Enum.AllowanceTypeSeller,
		Identifier: vo.Identifier,
		Buyer: vo.Buyer,
		Details: string(jsonStr),
		TaxAmount: taxAmount,
		TotalAmount: totalAmount,
		AllowanceStatus: Enum.AllowanceStatusInit,
	}
	if err := InvoiceDao.InsertAllowanceData(engine, data); err != nil {
		log.Error("Create Allowance Data Error", err)
		return err
	}
	if err := createAllowance(data); err != nil {
		log.Error("create Allowance Xml Error", err)
		return err
	}
	return nil
}
//?????????????????????
func cancelAllowance(data entity.AllowanceData) error {
	config := viper.GetStringMapString("PLATFORM")
	var resp InvoiceXml.CancelAllowance
	resp.CancelAllowanceNumber = data.AllowanceId
	resp.AllowanceDate = data.AllowanceDate
	resp.SellerId = config["company_ban"]
	resp.BuyerId = data.Identifier
	resp.CancelDate = data.CancelTime.Format("20060102")
	resp.CancelTime = data.CancelTime.Format("15:04:05")
	resp.CancelReason = data.CancelReason
	if err := GenerateCancelAllowanceXml(resp); err != nil {
		return err
	}
	return nil
}
//???????????????
func createAllowance(data entity.AllowanceData) error {
	var resp InvoiceXml.Allowance
	resp.Main = CreateAllowanceMain(data)
	resp.Details = CreateAllowanceDetails(data)
	resp.Amount.TotalAmount = data.TotalAmount
	resp.Amount.TaxAmount = data.TaxAmount
	if err := GenerateAllowanceXml(resp); err != nil {
		return err
	}
	return nil
}
//???????????????MAIN
func CreateAllowanceMain(data entity.AllowanceData) InvoiceXml.AllowanceMain {
	var resp InvoiceXml.AllowanceMain
	resp.AllowanceNumber = data.AllowanceId
	resp.AllowanceDate = data.AllowanceDate
	resp.Seller = generateInvoiceSeller()
	resp.Buyer = CreateAllowanceBuyer(data)
	resp.AllowanceType = data.AllowanceType
	return resp
}
//???????????????Details
func CreateAllowanceDetails(data entity.AllowanceData) InvoiceXml.AllowanceDetails {
	var resp InvoiceXml.AllowanceDetails
	var details []InvoiceVo.ProductItemVo
	if err := json.Unmarshal([]byte(data.Details), &details); err != nil {
		log.Error("Unmarshal Error", err)
	}
	for _, v := range details {
		var res InvoiceXml.AllowanceProductItem
		res.OriginalInvoiceNumber = v.OriginalInvoiceNumber
		res.OriginalInvoiceDate = v.OriginalInvoiceDate
		res.OriginalDescription = v.OriginalDescription
		res.OriginalSequenceNumber = v.OriginalSequenceNumber
		res.AllowanceSequenceNumber = v.AllowanceSequenceNumber
		res.Quantity = v.Quantity
		res.Amount = v.Amount
		res.Tax = v.Tax
		res.TaxType = v.TaxType
		res.Unit = v.Unit
		res.UnitPrice = v.UnitPrice
		resp.ProductItem = append(resp.ProductItem, res)
	}
	return resp
}
//???????????????Buyer??????
func CreateAllowanceBuyer(data entity.AllowanceData) InvoiceXml.InvoiceBuyer {
	var resp InvoiceXml.InvoiceBuyer
	resp.Identifier = data.Identifier
	resp.Name = data.Buyer
	return resp
}
//????????????
func VoidInvoice(data entity.InvoiceData) error {
	config := viper.GetStringMapString("PLATFORM")
	var resp InvoiceXml.VoidInvoiceC0701
	resp.VoidInvoiceNumber = fmt.Sprintf("%s%s", data.InvoiceTrack, data.InvoiceNumber)
	resp.InvoiceDate = data.CreateTime.Format("20060102")
	resp.SellerId = config["company_ban"]
	resp.BuyerId = data.Identifier
	resp.VoidDate = data.VoidTime.Format("20060102")
	resp.VoidTime = data.VoidTime.Format("15:04:05")
	resp.VoidReason = data.VoidReason
	if err := GenerateVoidInvoiceXml(resp); err != nil {
		return err
	}
	return nil
}
//????????????
func CancelInvoice(data entity.InvoiceData) error {
	config := viper.GetStringMapString("PLATFORM")
	var resp InvoiceXml.CancelInvoiceC0501
	resp.CancelInvoiceNumber = fmt.Sprintf("%s%s", data.InvoiceTrack, data.InvoiceNumber)
	resp.InvoiceDate = data.CreateTime.Format("20060102")
	resp.BuyerId = data.Identifier
	resp.SellerId = config["company_ban"]
	resp.CancelDate = data.CancelTime.Format("20060102")
	resp.CancelTime = data.CancelTime.Format("15:04:05")
	resp.CancelReason = data.CancelReason
	if err := GenerateCancelInvoiceXml(resp); err != nil {
		return err
	}
	return nil
}
//????????????
func CreateInvoice(engine *database.MysqlSession, vo InvoiceVo.InvoiceVo) (entity.InvoiceData, error) {
	var data entity.InvoiceData
	invoice, err := InvoiceDao.GetInvoiceByOrderId(engine, vo.OrderId)
	if err != nil {
		return data, err
	}
	if len(invoice.InvoiceNumber) == 0 {
		data, err = CreateInvoiceDataByOrder(engine, vo)
		if err != nil {
			return data, err
		}
	} else {
		data = invoice
	}
	var resp InvoiceXml.InvoiceC0401
	resp.Main = generateInvoiceMain(data)
	resp.Details.ProductItem = generateInvoiceDetails(data)
	resp.Amount = generateInvoiceAmount(data)
	if err := GenerateInvoiceXml(resp); err != nil {
		return data, err
	}
	return data, nil
}
//????????????MAIN
func generateInvoiceMain(data entity.InvoiceData) InvoiceXml.InvoiceC0401Main {
	var resp InvoiceXml.InvoiceC0401Main
	resp.InvoiceNumber = fmt.Sprintf("%s%s", data.InvoiceTrack, data.InvoiceNumber)
	resp.InvoiceDate = data.CreateTime.Format("20060102")
	resp.InvoiceTime = data.CreateTime.Format("15:04:05")
	resp.InvoiceType = data.InvoiceType
	resp.DonateMark = data.DonateMark
	if data.DonateMark == 1 {
		resp.NPOBAN = data.DonateBan
	}
	resp.PrintMark = data.PrintMark
	resp.CarrierType = Enum.CarrierType[data.CarrierType]
	resp.CarrierId1 = data.Carrier
	resp.CarrierId2 = data.Carrier
	resp.RandomNumber = data.RandomNumber
	resp.Seller = generateInvoiceSeller()
	resp.Buyer = generateInvoiceBuyer(data)
	resp.BuyerRemark = 1
	resp.CustomsClearanceMark = 1

	return resp
}
//????????????Seller??????
func generateInvoiceSeller() InvoiceXml.InvoiceSeller {
	config := viper.GetStringMapString("PLATFORM")
	var resp InvoiceXml.InvoiceSeller
	resp.Identifier = config["company_ban"]
	resp.Name = config["company_name"]
	resp.Address = config["company_address"]
	return resp
}
//????????????Buyer??????
func generateInvoiceBuyer(data entity.InvoiceData) InvoiceXml.InvoiceBuyer {
	var resp InvoiceXml.InvoiceBuyer
	resp.Identifier = data.Identifier
	resp.Name = data.Buyer
	return resp
}
//????????????Details??????
func generateInvoiceDetails(data entity.InvoiceData) []InvoiceXml.ProductItem {
	var resp []InvoiceXml.ProductItem

	var details []InvoiceVo.Details
	if err := json.Unmarshal([]byte(data.Detail), &details); err != nil {
		log.Error("Unmarshal Error", err)
	}
	for _, v := range details {
		var res InvoiceXml.ProductItem
		res.Description = v.ProductName
		res.Quantity = v.Quantity
		res.UnitPrice = v.ProductPrice
		res.Amount = v.ProductAmount
		res.SequenceNumber = v.Sequence
		resp = append(resp, res)
	}
	return resp
}
//????????????Amount??????
func generateInvoiceAmount(data entity.InvoiceData) InvoiceXml.InvoiceC0401Amount {
	var resp InvoiceXml.InvoiceC0401Amount
	resp.SalesAmount = float64(data.Sales)
	resp.FreeTaxSalesAmount = 0
	resp.ZeroTaxSalesAmount = 0
	resp.TaxType = Enum.InvoiceTaxType01
	resp.TaxRate = Enum.InvoiceTaxRate
	resp.TaxAmount = float64(data.Tax)
	resp.TotalAmount = data.Amount
	resp.DiscountAmount = 0
	resp.ExchangeRate = 0
	resp.OriginalCurrencyAmount = 0
	return resp
}
//??????????????????
func GetNextInvoiceNumber(engine *database.MysqlSession) (InvoiceVo.InvoiceResponse, error) {
	var resp InvoiceVo.InvoiceResponse
	yearMonth := tools.GetInvoiceYearMonth()
	data, err := InvoiceDao.GetInvoiceAssignNoByYearMonth(engine, yearMonth)
	if err != nil {
		return resp, err
	}
	if data.AssignId == 0 {
		return resp, fmt.Errorf("???????????????")
	}
	resp.Track = data.InvoiceTrack
	resp.Type = data.InvoiceType
	resp.Number = fmt.Sprintf("%0*v", 8, data.InvoiceNowNo)
	resp.Year = yearMonth[0:3]
	resp.Month = yearMonth[3:]
	//?????????????????????????????? ??????300???EMAIL?????? fixme ?????????????????? ???????????????
	last := tools.StringToInt64(data.InvoiceEndNo) - tools.StringToInt64(resp.Number)
	if  last < 300 {
		if last % 10 == 0 {
			if err := Mail.SendInvoiceSystemMail(last); err != nil {
				log.Error("Send Mail Error")
			}
		}
	}
	if resp.Number != data.InvoiceEndNo {
		data.InvoiceNowNo += 1
		if err := InvoiceDao.UpdateInvoiceAssignNoData(engine, data); err != nil {
			return resp, err
		}
	} else {
		data.InvoiceStatus = Enum.InvoiceAssignStatusDeEnabled
		if err := InvoiceDao.UpdateInvoiceAssignNoData(engine, data); err != nil {
			return resp, err
		}
	}
	return resp, nil
}
//??????????????????
func GetInvoiceAssignNumber() error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		return err
	}
	path := tools.GetFilePath("/invoice/E0501/", "", 0)
	result, err := GetInvoiceAssignNumberFile(path)
	if err != nil {
		log.Error("Get Invoice Assign Error", err)
		return err
	}
	for _, f := range result {
		var resp InvoiceXml.InvoiceAssignNo
		file, err := ioutil.ReadFile(path + f.Name())
		if err != nil {
			engine.Session.Rollback()
			log.Error("Open File Error", err)
			return err
		}
		value := xml.InvoiceXmlDecoder(string(file), "InvoiceEnvelope.InvoicePack.InvoiceAssignNo")
		if len(value.String()) != 0 {
			err := json.Unmarshal([]byte(value.String()), &resp)
			if err != nil {
				engine.Session.Rollback()
				log.Error("Json Unmarshal Error", err)
				return err
			}
		}
		if err := CreationInvoiceAssign(engine, resp); err != nil {
			engine.Session.Rollback()
			log.Error("Creation Invoice Assign Error", err)
			return err
		}
	}
	if err := engine.Session.Commit();err != nil {
		return err
	}
	return nil
}

func CreationInvoiceAssign(engine *database.MysqlSession, data InvoiceXml.InvoiceAssignNo) error {
	invoice, err := InvoiceDao.GetInvoiceAssignNoByYearMonthAndBooklet(engine, data.YearMonth, data.InvoiceBeginNo, data.InvoiceEndNo, data.InvoiceTrack)
	if err != nil {
		log.Error("Get Invoice Assign No Error", err)
		return err
	}
	if len(invoice.MonthYear) == 0 {
		var ent entity.InvoiceAssignNoData
		ent.InvoiceBan = data.Ban
		ent.InvoiceType = data.InvoiceType
		ent.MonthYear = data.YearMonth
		ent.InvoiceTrack = data.InvoiceTrack
		ent.InvoiceBeginNo = data.InvoiceBeginNo
		ent.InvoiceEndNo = data.InvoiceEndNo
		booklet, err := strconv.Atoi(data.InvoiceBeginNo)
		if err != nil {
			return err
		}
		ent.InvoiceBooklet = int64(booklet)
		beginNo, err := strconv.Atoi(data.InvoiceBeginNo)
		if err != nil {
			return err
		}
		ent.InvoiceNowNo = int64(beginNo)
		ent.InvoiceStatus = Enum.InvoiceAssignStatusDeEnabled
		if err := InvoiceDao.InsertInvoiceAssignNoData(engine, ent); err != nil {
			log.Error("Insert Invoice Assign No Error", err)
			engine.Session.Rollback()
			return err
		}
	}
	return nil
}

//????????????????????????
func CreateInvoiceDataByOrder(engine *database.MysqlSession, vo InvoiceVo.InvoiceVo) (entity.InvoiceData, error) {
	var ent entity.InvoiceData
	data, err := GetNextInvoiceNumber(engine)
	if err != nil {
		return ent, err
	}
	ent.InvoiceNumber = data.Number
	ent.InvoiceTrack = data.Track
	ent.InvoiceType = data.Type
	ent.Month = data.Month
	ent.Year = data.Year
	ent.OrderId = vo.OrderId
	ent.BuyerId = vo.UserId
	jsonStr, _ := json.Marshal(vo.Details)
	ent.Detail = string(jsonStr)
	ent.Amount = vo.Amount
	ent.RandomNumber = tools.RangeNumber(9999, 4)
	if vo.Carrier.InvoiceType != Enum.InvoiceTypeCompany {
		//??????????????? ???????????????10???0
		ent.Identifier = "0000000000"
		ent.Buyer = vo.UserId
	} else {
		ent.Identifier = vo.Carrier.CompanyBan
		ent.Buyer = vo.Carrier.CompanyName
	}
	//???????????? 0???????????? 1:??????(???????????????????????????)
	if vo.Carrier.InvoiceType == Enum.InvoiceTypeDonate {
		ent.DonateMark = 1
		ent.DonateBan = vo.Carrier.DonateBan
	} else {
		ent.DonateMark = 0
	}
	ent.CarrierType = vo.Carrier.CarrierType
	ent.Carrier = vo.Carrier.CarrierId
	ent.PrintMark = "N"
	ent.Tax = int64(tools.Round(float64(vo.Amount) - (float64(vo.Amount)/Enum.InvoiceTaxRate)))
	ent.Sales = vo.Amount - ent.Tax
	ent.InvoiceStatus = Enum.InvoiceStatusNot
	ent.CreateTime = time.Now()
	log.Debug("ent", ent)
	if err := InvoiceDao.InsertInvoiceData(engine, ent); err != nil {
		return ent, err
	}
	return ent, nil
}

func GetDonateInfo(engine *database.MysqlSession, DonateCode string) string {
	donate, err := member.GetDonateDataByDonateCode(engine, DonateCode)
	if err != nil {
		log.Error("Get Donate Unit Data Error", err)
		return ""
	}
	return fmt.Sprintf("%s/????????????%s", donate.DonateShort, donate.DonateCode)
}
//??????????????????
func ReadAwardedPathFiles() error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	path := fmt.Sprintf("./data/Invoice/awarded/src/")
	result, err := ioutil.ReadDir(path)
	if err != nil {
		log.Error("Open File Error", err)
		return err
	}
	for _, f := range result {
		log.Debug("files", f.Name())
		data, err := ReadAwardedFile(f.Name())
		if err != nil {
			log.Error("Error when opening file:", err)
			return err
		}
		log.Debug("sss", data)
		if err := processAwarded(engine, data, f.Name()[0:1]); err != nil {
			return err
		}
		if err := tools.MoveAwardedFile(f.Name()); err != nil {
			return err
		}
	}
	return nil
}
//??????????????????
func ReadAwardedFile(filename string) ([]InvoiceVo.Awarded, error){
	path := fmt.Sprintf("./data/Invoice/awarded/src/%s", filename)
	file, err := os.Open(path)
	if err != nil {
		log.Error("Error when opening file:", err)
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	var data []InvoiceVo.Awarded
	//????????????
	for fileScanner.Scan() {
		body := InvoiceVo.Awarded{}
		NewInvoiceRule().SetInvoiceData(fileScanner.Text(), &body)
		data = append(data, body)
	}
	if err := fileScanner.Err(); err != nil {
		log.Error("Error While Reading File:", err)
		return nil, err
	}
	file.Close()
	return data, nil
}
//??????????????????
func processAwarded(engine *database.MysqlSession, data []InvoiceVo.Awarded, model string) error {
	var year, month string
	for _, v := range data {
		if !validate.IsVerifyEnglish(v.InvoiceAxle) {
			continue
		}
		result, err := InvoiceDao.GetInvoiceByTrackAndNumber(engine, v.InvoiceAxle, v.InvoiceNumber)
		if err != nil {
			log.Error("Get Invoice Data Error", err)
			return err
		}
		if result.InvoiceId == 0 {
			continue
		}
		if len(year) == 0 {
			year = result.Year
			month = result.Month
		}
		result.AwardModel = model
		result.AwardTime = time.Now()
		result.InvoiceStatus = Enum.InvoiceStatusWin
		if err := InvoiceDao.UpdateInvoiceData(engine, result); err != nil {
			log.Error("Update Invoice Data Error", err)
			return err
		}
		Month := tools.GetYearMonth(result.Month)
		//??????????????????
		message := fmt.Sprintf("?????????\n??????????????????????????????????????? Check???Ne ???????????????????????????(???????????????%s???%s???%s???)?????????????????????\n" +
			"?????????????????? Check???Ne ????????????????????????????????????????????????????????????????????? Check???Ne ??????????????????????????????????????????????????????????????????????????????????????????????????????\n" +
			"Check???Ne ????????????", result.Year, Month[0], Month[1])
		if err := Notification.SendSystemNotify(engine, result.BuyerId, message, Enum.NotifyMsgTypePlaPlatform, ""); err != nil {
			return err
		}
	}
	if err := InvoiceDao.UpdateInvoiceAllData(engine, year, month); err != nil {
		log.Error("Update Invoice Data Error", err)
		return err
	}
	return nil
}

//????????????
func ResendInvoice(params Request.ErpSearchOrderRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if len(params.OrderId) == 0 {
		return fmt.Errorf("1001007")
	}
	data, err := InvoiceDao.GetInvoiceByOrderId(engine, params.OrderId)
	if err != nil {
		return fmt.Errorf("1001007")
	}
	if len(data.OrderId) == 0 {
		log.Error("????????????")
		return fmt.Errorf("1001007")
	}
	var resp InvoiceXml.InvoiceC0401
	resp.Main = generateInvoiceMain(data)
	resp.Details.ProductItem = generateInvoiceDetails(data)
	resp.Amount = generateInvoiceAmount(data)
	if err := GenerateInvoiceXml(resp); err != nil {
		log.Error("Generate Invoice Xml Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}
//??????????????????
func HandleNewInvoiceTrack(params Request.InvoiceTrackRequest) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	if err := engine.Session.Begin(); err != nil {
		log.Error("Commit Begin Invoice Assign No Error", err)
		return fmt.Errorf("1001001")
	}
	var ent entity.InvoiceAssignNoData
	ent.InvoiceBan = viper.GetString("PLATFORM.COMPANY_BAN")
	ent.InvoiceType = "7"
	ent.MonthYear = params.InvoicePeriod
	ent.InvoiceTrack = params.InvoiceTrack
	ent.InvoiceBeginNo = params.InvoiceBegin
	ent.InvoiceEndNo = params.InvoiceEnd
	booklet := tools.StringToInt64(params.InvoiceEnd)
	ent.InvoiceBooklet = booklet + 1
	beginNo := tools.StringToInt64(params.InvoiceBegin)
	ent.InvoiceNowNo = beginNo
	ent.InvoiceStatus = Enum.InvoiceAssignStatusEnable
	if err := InvoiceDao.InsertInvoiceAssignNoData(engine, ent); err != nil {
		log.Error("Insert Invoice Assign No Error", err)
		engine.Session.Rollback()
		return fmt.Errorf("1001001")
	}
	if err := engine.Session.Commit();err != nil {
		log.Error("Commit Invoice Assign No Error", err)
		return fmt.Errorf("1001001")
	}
	return nil
}
//
//func generateTrackBlank()  {
//
//}