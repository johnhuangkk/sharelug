package controllers

import (
	"api/services/Enum"
	"api/services/Service/CvsShipping"
	"api/services/Service/SevenMyshipApi"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"api/services/util/response"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Countries struct {
	City []string `json:"city"`
}
type DistrictList struct {
	District []string
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
func GetCityShopsAddress(ctx *gin.Context) {

	res := response.New(ctx)
	country := ctx.Query("country")
	district := ctx.Query("district")
	shopData, _ := SevenMyshipApi.GetShopAddress(country, district)
	if shopData == nil {
		res.Success("OK").SetData(make([]string, 0)).Send()
	} else {

		res.Success("OK").SetData(shopData).Send()
	}
}

// CreateOLTP 建立OLTP
func CreateOLTP(ctx *gin.Context) {

	body, err := ioutil.ReadAll(ctx.Request.Body)

	// defer ctx.Request.Body.Close()
	if err != nil {
		log.Error("Seven OLTP error", err.Error())
		return
	}

	s, datas := SevenMyshipApi.CreateChargeOrderRecordByPos(body)

	ctx.Data(http.StatusOK, `application/xml`, s)

	engine := database.GetMysqlEngine()
	defer engine.Close()
	var UpdateCvsShipping CvsShipping.UpdateCvsShipping

	UpdateCvsShipping.ShipType = Enum.CVS_7_ELEVEN
	var shipMap entity.SevenShipMapData
	for _, data := range datas {
		if data.OLPrint == "Y" && (data.OLOiNo == "850" || data.OLOiNo == "851") {
			t, _ := time.Parse(`20060102150405`, data.Date+data.Time)

			flag, err := engine.Engine.Table("seven_ship_map_data").Select("paymentno_with_code,payment_no").Where("payment_no =?", data.OLCode2[0:8]).Get(&shipMap)
			if flag {
				UpdateCvsShipping.ShipNo = shipMap.PaymentNoWithCode
				UpdateCvsShipping.Type = "000"
				UpdateCvsShipping.DateTime = t.Format(`2006-01-02 15:04:05`)
				UpdateCvsShipping.DetailStatus = data.OLPrint
				UpdateCvsShipping.FlowType = "N"
				UpdateCvsShipping.Log = string(s)
				UpdateCvsShipping.FileName = ""
				UpdateCvsShipping.OrderNo = shipMap.OrderId
				err = UpdateCvsShipping.UpdateCvsShippingShipment(engine)

				if err != nil {
					log.Error("Seven OLTP error", err)
					return
				}
			} else {
				log.Error("Seven OLTP error", err)
				return
			}
			log.Info("Seven API寄取即時更新", t.Format(`2006-01-02 15:04:05`), shipMap.PaymentNoWithCode)
		} else if data.OLPrint == "Y" && data.OLOiNo == "852" {
			t, _ := time.Parse(`20060102150405`, data.Date+data.Time)
			flag, err := engine.Engine.Table("seven_ship_map_data").Select("paymentno_with_code,payment_no,ship_no").Where("ship_no =?", data.OLCode2[0:8]).Get(&shipMap)
			if flag {
				UpdateCvsShipping.ShipNo = shipMap.PaymentNoWithCode
				UpdateCvsShipping.Type = "777"
				UpdateCvsShipping.DateTime = t.Format(`2006-01-02 15:04:05`)
				UpdateCvsShipping.DetailStatus = data.OLPrint
				UpdateCvsShipping.FlowType = "N"
				UpdateCvsShipping.Log = string(s)
				UpdateCvsShipping.FileName = ""
				UpdateCvsShipping.OrderNo = shipMap.OrderId
				err = UpdateCvsShipping.UpdateCvsShippingSuccess(engine)

				if err != nil {
					log.Error("Seven OLTP error", err)
					return
				}
			} else {
				log.Error("Seven OLTP error", err)
				return
			}
			log.Info("Seven API寄取即時更新", t.Format(`2006-01-02 15:04:05`), shipMap.PaymentNoWithCode)
		} else if data.OLPrint == "Y" && data.OLOiNo == "853" {
			t, _ := time.Parse(`20060102150405`, data.Date+data.Time)

			flag, err := engine.Engine.Table("seven_ship_map_data").Select("paymentno_with_code,payment_no").Where("payment_no =?", data.OLCode2[0:8]).Get(&shipMap)
			if flag {
				UpdateCvsShipping.ShipNo = shipMap.PaymentNoWithCode
				UpdateCvsShipping.Type = "888"
				UpdateCvsShipping.DateTime = t.Format(`2006-01-02 15:04:05`)
				UpdateCvsShipping.DetailStatus = data.OLPrint
				UpdateCvsShipping.FlowType = "R"
				UpdateCvsShipping.Log = string(s)
				UpdateCvsShipping.FileName = ""
				UpdateCvsShipping.OrderNo = shipMap.OrderId
				flag, err := checkCvsExist(engine, shipMap.PaymentNoWithCode, UpdateCvsShipping.Type)
				if !flag {
					err = UpdateCvsShipping.OnlyWriteShippingLog(engine, true)
				}

				if err != nil {
					log.Error("Seven OLTP error", err)
					return
				}
			} else {
				log.Error("Seven OLTP error", err)
				return
			}
			log.Info("Seven API寄取即時更新", t.Format(`2006-01-02 15:04:05`), shipMap.PaymentNoWithCode)
		}

	}

}
func checkCvsExist(engine *database.MysqlSession, paymentNoWithCode string, storeType string) (bool, error) {
	var cvsShippingLogData entity.CvsShippingLogData
	flag, err := engine.Engine.Table("cvs_shipping_log_data").Select("ship_no,cvs_type,type").Where("ship_no =? && cvs_type =? && type=?", paymentNoWithCode, Enum.CVS_7_ELEVEN, storeType).Get(&cvsShippingLogData)
	return flag, err
}
func FetchDailyShops(ctx *gin.Context) {
	res := response.New(ctx)
	ctx.SetCookie("site_cookie", "cookievalue", 3600, "/", "localhost", false, true)

	SevenMyshipApi.FetchDailyShopStatus()
	res.Success("OK").Send()
}

func Test(ctx *gin.Context) {
	res := response.New(ctx)
	// SevenMyshipApi.FetchPackageSendByOl()
	SevenMyshipApi.FetchCPPSStatus()
	// SevenMyshipApi.FetchCEINStatus()
	// SevenMyshipApi.FetchCEDRStatus()
	// SevenMyshipApi.FetchCERTStatus()
	// SevenMyshipApi.FetchCESPStatus()
	// PostBag.UploadShipOrderFile()

	res.Success("OK").Send()
}
