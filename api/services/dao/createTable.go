package main

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"bytes"
	"flag"
	"io/ioutil"

	"github.com/spf13/viper"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "config/config.yaml", "Configuration file path.")
	flag.Parse()

	err := readConfig()
	if err != nil {
		panic(err)
	}
}

//讀取config的檔案
func readConfig() error {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = viper.ReadConfig(bytes.NewBuffer(content))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	createTable(engine, entity.OrderData{})
	createTable(engine, entity.MemberData{})
	createTable(engine, entity.OrderDetail{})
	createTable(engine, entity.OtpData{})
	createTable(engine, entity.ProductData{})
	createTable(engine, entity.ProductSpecData{})
	createTable(engine, entity.StoreData{})
	createTable(engine, entity.ProductImagesData{})
	createTable(engine, entity.DeviceData{})
	createTable(engine, entity.TransferData{})
	createTable(engine, entity.TransferLogData{})
	createTable(engine, entity.GwCreditAuthData{})
	createTable(engine, entity.MemberCardData{})
	createTable(engine, entity.GwCreditAuthLog{})
	createTable(engine, entity.TwidCheckLogData{})
	createTable(engine, entity.PostBoxData{})
	createTable(engine, entity.PostShippingStatus{})
	createTable(engine, entity.PostConsignmentData{})
	createTable(engine, entity.Sequence{})
	createTable(engine, entity.MartFamilyStoreData{})
	createTable(engine, entity.MartFamilyShippingData{})
	createTable(engine, entity.MartOkStoreData{})
	createTable(engine, entity.MartOkShippingData{})
	createTable(engine, entity.MartHiLifeStoreData{})
	createTable(engine, entity.MartHiLifeShippingData{})
	createTable(engine, entity.MartShippingFileData{})
	createTable(engine, entity.UserAddressData{})
	createTable(engine, entity.TwIdVerifyData{})
	createTable(engine, entity.ScsBankAccountData{})
	createTable(engine, entity.OrderMessageBoardData{})
	createTable(engine, entity.StatusHistoryLog{})
	createTable(engine, entity.StoreRankData{})
	createTable(engine, entity.BalanceAccountData{})
	createTable(engine, entity.BalanceRetainAccountData{})
	createTable(engine, entity.SystemLog{})
	createTable(engine, entity.PlatformOrderData{})
	createTable(engine, entity.OrderRefundData{})
	createTable(engine, entity.CreditBatchRequestData{})
	createTable(engine, entity.AccountActivityData{})
	createTable(engine, entity.NotifyMessageData{})
	createTable(engine, entity.WithdrawData{})
	createTable(engine, entity.BankCodeData{})
	createTable(engine, entity.EmailVerifyData{})
	createTable(engine, entity.MemberWithdrawData{})
	createTable(engine, entity.B2cOrderData{})
	createTable(engine, entity.B2cBillingData{})
	createTable(engine, entity.CvsShippingData{})
	createTable(engine, entity.CvsShippingLogData{})
	createTable(engine, entity.UpgradeProductData{})
	createTable(engine, entity.ContactData{})
	createTable(engine, entity.OnlineNotifyData{})
	createTable(engine, entity.SmsLogData{})
	createTable(engine, entity.CustomerQuestionData{})
	createTable(engine, entity.CustomerData{})

	createTable(engine, entity.SevenMyshipShopData{})
	createTable(engine, entity.SevenChargeOrderData{})
	createTable(engine, entity.BankBinCode{})
	createTable(engine, entity.CvsAccountingData{})
	createTable(engine, entity.SevenShipMapData{})
	createTable(engine, entity.InvoiceAssignNoData{})
	createTable(engine, entity.InvoiceData{})
	createTable(engine, entity.AllowanceData{})
	createTable(engine, entity.CancelAllowanceData{})
	createTable(engine, entity.MemberCarrierData{})
	createTable(engine, entity.DonateData{})
	createTable(engine, entity.PostBagConsignmentData{})
	createTable(engine, entity.StoreSocialMediaData{})
	createTable(engine, entity.ShortUrlData{})
	createTable(engine, entity.BillOrderData{})
	createTable(engine, entity.ErpUser{})
	createTable(engine, entity.CustomerMemo{})
	createTable(engine, entity.BatchShipExcelImport{})
	createTable(engine, entity.TaiwanArea{})
	createTable(engine, entity.StoreSelfDeliveryArea{})

	createTable(engine, entity.Promotion{})
	createTable(engine, entity.PromotionCode{})
	createTable(engine, entity.CouponUsedRecord{})
	createTable(engine, entity.PromoCodeOperateRecord{})
	createTable(engine, entity.KgiSpecialStore{})
	createTable(engine, entity.IndustryData{})
	createTable(engine, entity.TaiwanCity{})
	createTable(engine, entity.MemberSendKgiBank{})
	createTable(engine, entity.ProductHistoryLog{})
	createSequenceData(engine)
}

func createTable(engine *database.MysqlSession, beanOrTableName interface{}) {

	isExist, err := engine.Session.IsTableExist(beanOrTableName)
	if err != nil {
		log.Debug("is table exist error : %v", err)
		panic(err)
	}
	if !isExist {
		if err := engine.Session.CreateTable(beanOrTableName); err != nil {
			log.Debug("create table error : %v", err)
			panic(err)
		}
	} else {
		_ = engine.Session.Sync2(beanOrTableName)
	}
}

func createSequenceData(engine *database.MysqlSession) {
	sql := "INSERT INTO sequence (name, current_value, increment) VALUES (?, ?, ?)"

	if !getSequence(engine, "KgiAtmBank") {
		_, err := engine.Session.Exec(sql, "KgiAtmBank", 1, 1)
		if err != nil {
			log.Debug("insert data error!", err)
		}
	}
	if !getSequence(engine, "qrcode") {
		_, err := engine.Session.Exec(sql, "qrcode", 1, 1)
		if err != nil {
			log.Debug("insert data error!", err)
		}
	}
	if !getSequence(engine, "iPost") {
		_, err := engine.Session.Exec(sql, "iPost", 1, 1)
		if err != nil {
			log.Debug("insert data error!", err)
		}
	}
	if !getSequence(engine, "order") {
		_, err := engine.Session.Exec(sql, "order", 1, 1)
		if err != nil {
			log.Debug("insert data error!", err)
		}
	}
	if !getSequence(engine, "postBag") {
		_, err := engine.Session.Exec(sql, "postBag", 1, 1)
		if err != nil {
			log.Debug("insert postBag sequence error!", err)
		}
	}
	if !getSequence(engine, "customerQuestion") {
		_, err := engine.Session.Exec(sql, "customerQuestion", 1, 1)
		if err != nil {
			log.Debug("insert customerQuestion sequence error!", err)
		}
	}
	if !getSequence(engine, "taiwanArea") {
		_, err := engine.Session.Exec(sql, "taiwanArea", 1, 1)
		if err != nil {
			log.Debug("insert taiwanArea sequence error!", err)
		}
	}
	if !getSequence(engine, "couponAction") {
		_, err := engine.Session.Exec(sql, "couponAction", 1, 1)
		if err != nil {
			log.Debug("insert taiwanArea sequence error!", err)
		}
	}
}

func getSequence(engine *database.MysqlSession, key string) bool {
	count, err := engine.Session.Table("sequence").
		Where("name = ?", key).Count()
	if err != nil {
		panic(err)
	}
	return count != 0
}
