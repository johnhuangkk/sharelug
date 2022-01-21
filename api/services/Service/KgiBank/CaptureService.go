package KgiBank

import (
	"api/services/Enum"
	"api/services/Service/Mail"
	"api/services/Service/OrderService"
	"api/services/VO/KgiBank"
	"api/services/dao/Credit"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	tools "api/services/util/tools"
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func HandleC2C3DCapture() {
	var MerchantID = viper.GetString("KgiCredit.C2C.3D.MerchantID")
	var TerminalID = viper.GetString("KgiCredit.C2C.3D.TerminalID")
	data, total, _ := GetCapture(MerchantID, TerminalID, Enum.OrderTransC2c, Enum.OrderTrans3D)
	log.Debug("count", len(data), total)
	err := ExporterKgiFile(data, MerchantID, total)
	if err != nil {
		log.Error("Exporter Kgi File")
	}
}
func HandleC2CN3DCapture() {
	var MerchantID = viper.GetString("KgiCredit.C2C.N3D.MerchantID")
	var TerminalID = viper.GetString("KgiCredit.C2C.N3D.TerminalID")
	data, total, _ := GetCapture(MerchantID, TerminalID, Enum.OrderTransC2c, Enum.OrderTransN3D)
	log.Debug("count", len(data), total)
	err := ExporterKgiFile(data, MerchantID, total)
	if err != nil {
		log.Error("Exporter Kgi File")
	}
}

func HandleB2C3DCapture() {
	var MerchantID = viper.GetString("KgiCredit.B2C.3D.MerchantID")
	var TerminalID = viper.GetString("KgiCredit.B2C.3D.TerminalID")
	data, total, _ := GetCapture(MerchantID, TerminalID, Enum.OrderTransB2c, Enum.OrderTrans3D)
	log.Debug("count", len(data), total)
	err := ExporterKgiFile(data, MerchantID, total)
	if err != nil {
		log.Error("Exporter Kgi File")
	}
}
func HandleB2CN3DCapture() {
	var MerchantID = viper.GetString("KgiCredit.B2C.N3D.MerchantID")
	var TerminalID = viper.GetString("KgiCredit.B2C.N3D.TerminalID")
	data, total, _ := GetCapture(MerchantID, TerminalID, Enum.OrderTransB2c, Enum.OrderTransN3D)
	log.Debug("count", len(data), total)
	err := ExporterKgiFile(data, MerchantID, total)
	if err != nil {
		log.Error("Exporter Kgi File")
	}
}


func GetCapture(MerchantID, TerminalID, payType, creditType string) ([]KgiBank.Body, int, error)  {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	var body []KgiBank.Body
	total := 0
	data, err := Credit.GetGwCreditBySuccess(engine, payType, creditType)
	if err != nil {
		return body, total, err
	}

	for _, v := range data {
		var str KgiBank.Body
		str.MerchantId = v.MerchantId
		str.TerminalId = v.TerminalId
		str.OrderId = v.OrderId      //訂單號碼
		str.TranAmount = strconv.FormatInt(v.TramsAmount, 10)   //交易金額
		str.AuthCode = v.ApproveCode    //授權碼
		str.TranDate = v.TransTime.Format("20060102")     //交易日期
		str.ProcessDate = time.Now().Format("060102")  //帳單處理日期
		str.ResponseCode = v.ResponseCode //回應碼
		str.ResponseMsg = v.ResponseMsg  //回應訊息
		//str.PaymentDate =
		if v.TransType == Enum.CreditTransTypeRefund {
			str.TranType = "01"     //交易碼
			total -= int(v.TramsAmount)
		} else {
			str.TranType = "02"     //交易碼
			total += int(v.TramsAmount)
		}
		body = append(body, str)
	}
	return body, total, nil
}

//產生請款檔
func ExporterKgiFile(data []KgiBank.Body, MerchantID string, ToTal int) error {

	engine := database.GetMysqlEngine()
	defer engine.Close()

	var ent entity.CreditBatchRequestData
	ent.SendDate = time.Now().Format("20060102")
	ent.Count = strconv.Itoa(len(data))
	ent.Symbol = tools.Sign(float64(ToTal))
	ent.ToTalAmount = strconv.Itoa(int(math.Abs(float64(ToTal))))
	ent.Status = "INIT"
	//寫入批次記錄
	Batch, err := Credit.InsertCreditBatchRequestData(engine, ent)
	if err != nil {
		log.Debug("Writer File Error", err)
		return err
	}
	var fileContent string
	header := KgiBank.Header{
		Flag: "H",
		MerchantId: MerchantID,
		SendDate: time.Now().Format("20060102"),
		Seq: strconv.Itoa(Batch.BatchId),
		Count: Batch.Count,
		Symbol: Batch.Symbol,
		ToTal: Batch.ToTalAmount,
		Filler: "",
	}
	content := NewKgiHeaderRule().ToString(header)
	fileContent += content
	if len(data) == 0 {
		log.Info("上傳結束")
		return nil
	}
	for _, v := range data {
		content := NewKgiBodyRule().ToString(v)
		fileContent += content
		_ = Credit.UpdateGwCreditBatchId(engine, v.OrderId, v.AuthCode, strconv.Itoa(Batch.BatchId))
	}
	path := tools.GetFilePath("/kgi/export/", "", 0)
	log.Debug("data len", path, len(fileContent))
	fileName := fmt.Sprintf("%s.dat", MerchantID)
	_, err = tools.CreateFile(path, fileContent, fileName)
	if err != nil {
		log.Debug("Writer File Error", err)
	}
	err = uploadFile(MerchantID)
	if err != nil {
		return err
	}

	return nil
}


//上傳檔案到銀行
func uploadFile(MerchantID string) error {
	var hostname = viper.GetString("KgiCredit.hostname")
	cmd := exec.Command("java", "-jar", "MDECFileUploadTls1.2.jar", "Internet", MerchantID, hostname)
	cmd.Dir = tools.GetFilePath("/kgi/export/", "", 0)
	out, err := cmd.Output()
	if err != nil {
		log.Error("[kgi] Upload File Error", err)
		//return err
	}
	filename := fmt.Sprintf("%s.dat", time.Now().Format("20060102150405"))
	log.Error("file", MerchantID)
	path := tools.GetFilePath("/kgi/export/", "", 0)
	if err := os.Rename(path + MerchantID + ".dat", tools.GetFilePath("/kgi/export/succ/", "", 0) + filename); err != nil {
		log.Error("Rename Error", err)
		return err
	}
	log.Info("[kgi] Upload File", tools.Big5ToUtf8(string(out)))
	return nil
}

//下載請款檔
func DownloadFile(date time.Time, MerchantID string) error {
	var hostname = viper.GetString("KgiCredit.hostname")
	cmd := exec.Command("java", "-jar", "MDECFileDownloadTls1.2.jar", "Internet", MerchantID, hostname, date.Format("20060102"))
	path := tools.GetFilePath("/kgi/import/", "", 0)
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("[kgi] Download File Error", err)
		return err
	}
	log.Info("[kgi] Upload File", string(out))
	if err := HandleRespond(); err != nil {
		log.Error("Handle Respond Error", err)
		return nil
	}
	return nil
}
//信用卡請款回檔處理
func HandleRespond() error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	//
	folder := "./data/kgi/import/"
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		log.Debug("folder files :", file.Name())
		if file.IsDir() || file.Name() == "MDECFileDownloadTls1.2.jar" || file.Name() == "recv" || file.Name() == "succ" {
			continue
		} else {
			data, err := ReadFile(file.Name())
			if err != nil {
				return err
			}
			for _, v := range data {
				log.Debug("data row", v.OrderId, v.ResponseCode, v.ResponseMsg, v.TranType)
				TranType := Enum.CreditTransTypeAuth
				if v.TranType == "01" {
					TranType = Enum.CreditTransTypeRefund
				}
				GwData, err := Credit.GetGwCreditByOrderIdAndTranType(engine, v.OrderId, TranType)
				if err != nil {
					return err
				}
				if GwData.CaptureStatus != Enum.CreditCaptureInit {
					continue
				}
				if v.ResponseCode == "00" {
					GwData.CaptureStatus = Enum.CreditCaptureSuccess
				} else {
					GwData.CaptureStatus = Enum.CreditCaptureFail
					if v.ResponseCode == "A05" {
						if err := Mail.SendCreditSystemMail(v.OrderId); err != nil {
							log.Error("Send Mail Error")
						}
					}
				}
				GwData.CaptureCode = v.ResponseCode
				GwData.CaptureMsg = v.ResponseMsg
				GwData.CaptureTime = time.Now()
				if err := Credit.UpdateGwCreditData(engine, GwData); err != nil {
					return err
				}
				if GwData.TransType == Enum.CreditTransTypeRefund && GwData.CaptureStatus == Enum.CreditCaptureSuccess {
					if err := OrderService.ChangeReturnStatus(engine, GwData.OrderId); err != nil {
						return err
					}
				}
			}
			if err := os.Rename("./data/kgi/import/" + file.Name(), "./data/kgi/import/succ/" + file.Name()); err != nil {
				log.Error("Rename Error", err)
				return err
			}
		}
	}
	return nil
}

func ReadFile(filename string) ([]KgiBank.Body, error){
	path := fmt.Sprintf("./data/kgi/import/%s", filename)
	file, err := os.Open(path)
	if err != nil {
		log.Error("Error when opening file:", err)
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)
	var data []KgiBank.Body
	//讀取單行
	for fileScanner.Scan() {
		body := KgiBank.Body{}
		charset := tools.Big5ToUtf8(fileScanner.Text())
		NewKgiBodyRule().SetRawData(charset, &body)
		data = append(data, body)
	}
	if err := fileScanner.Err(); err != nil {
		log.Error("Error While Reading File:", err)
		return nil, err
	}
	file.Close()
	return data, nil
}