package tools

import (
	"api/services/dao/sequence"
	"api/services/util/log"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

//產生 Order Id
func GeneratorOrderId() string {
	id, err := sequence.GetOrderIdSeq()
	if err != nil {
		log.Error("Sequence Database Error", err)
	}
	content := fmt.Sprintf("%0*s", 4, id)
	return fmt.Sprintf("B%s%s%s", time.Now().Format("060102"), content[len(content)-4:], RangeNumber(99, 2))
}

func GeneratorMarketingOrderId() string {
	id, err := sequence.GetOrderIdSeq()
	if err != nil {
		log.Error("Sequence Database Error", err)
	}
	content := fmt.Sprintf("%0*s", 4, id)
	return fmt.Sprintf("BM%s%s%s", time.Now().Format("060102"), content[len(content)-4:], RangeNumber(99, 2))
}


func GeneratorRealtimeOrderId() string {
	id, err := sequence.GetOrderIdSeq()
	if err != nil {
		log.Error("Sequence Database Error", err)
	}
	content := fmt.Sprintf("%0*s", 4, id)
	return fmt.Sprintf("R%s%s%s", time.Now().Format("060102"), content[len(content)-4:], RangeNumber(99, 2))
}

func GeneratorBillOrderId() string {
	id, err := sequence.GetOrderIdSeq()
	if err != nil {
		log.Error("Sequence Database Error", err)
	}
	content := fmt.Sprintf("%0*s", 4, id)
	return fmt.Sprintf("A%s%s%s", time.Now().Format("060102"), content[len(content)-4:], RangeNumber(99, 2))
}

//產生 Product Id
func GeneratorProductId() string {
	return fmt.Sprintf("I%s%s", time.Now().Format("0601021504"), RangeNumber(99999, 5))
}

//產生 Product Spec Id
func GeneratorProductSpecId(productId string, number int) string {
	return fmt.Sprintf("%s-%03d", productId, number)
}

func GeneratorCardId() string {
	return fmt.Sprintf("MC%s%s", time.Now().Format("060102150405"), RangeNumber(99999, 5))
}

func GeneratorOrderRefundId() string {
	return fmt.Sprintf("RF%s%s", time.Now().Format("060102150405"), RangeNumber(999, 3))
}

func GeneratorOrderReturnId() string {
	return fmt.Sprintf("RT%s%s", time.Now().Format("060102150405"), RangeNumber(999, 3))
}

func GeneratorWithdrawId() string {
	return fmt.Sprintf("W%s%s", time.Now().Format("0601021504"), RangeNumber(999, 3))
}

func GeneratorTransId() string {
	return fmt.Sprintf("00%s%s", time.Now().Format("0102"), RangeNumber(99, 2))
}

//發票折讓單編號
func GeneratorAllowanceId() string {
	return fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), RangeNumber(99, 2))
}

//產生店家ID
func GeneratorStoreId() string {
	return fmt.Sprintf("S%s%s", time.Now().Format("060102150405"), RangeNumber(99999, 5))
}

func GeneratorValidationCode(code string) string {
	code += time.Now().Format("2006/01/02 15:04:05")
	Key := viper.GetString("EncryptKey")
	return fmt.Sprintf("%x", sha256.Sum256([]byte(AesEncrypt(code, Key))))
}

func GeneratorB2COrderId() string {
	id, err := sequence.GetOrderIdSeq()
	if err != nil {
		log.Error("Sequence Database Error", err)
	}
	content := fmt.Sprintf("%0*s", 4, id)
	return fmt.Sprintf("BC%s%s%s", time.Now().Format("060102"), content[len(content)-4:], RangeNumber(99, 2))
}

func GeneratorB2CBillId() string {
	id, err := sequence.GetOrderIdSeq()
	if err != nil {
		log.Error("Sequence Database Error", err)
	}
	content := fmt.Sprintf("%0*s", 4, id)
	return fmt.Sprintf("BILL%s%s%s", time.Now().Format("060102"), content[len(content)-4:], RangeNumber(99, 2))
}

func GeneratorBatchShipId() string {
	return fmt.Sprintf("%s%s", time.Now().Format("060102150405"), RangeNumber(99999, 5))
}

func GeneratorAchId() string {
	id, err := sequence.GetOrderIdSeq()
	if err != nil {
		log.Error("Sequence Database Error", err)
	}
	return fmt.Sprintf("%0*s", 6, id)
}

// gen customer question id
func GenerateCustomerId() string {
	now := time.Now().Format("20060102150405")
	id, err := sequence.GetCustomerQuestionSeq()
	if err != nil {
		log.Error("Get Customer Database Sequence", err)
	}
	content := fmt.Sprintf("%0*s", 4, id)
	return fmt.Sprintf("C%s%s", now, content)
}

func GenerateTaiwanAreaId() string {
	id, err := sequence.GetTaiwanAreaSeq()
	if err != nil {
		log.Error("Get Customer Database Sequence", err)
	}
	return fmt.Sprintf("T%s", id)
}
func GenerateCouponActionId() string {
	id, err := sequence.GetCouponActionSeq()
	if err != nil {
		log.Error("Get CouponAction Sequence", err)
	}
	content := fmt.Sprintf("%0*s", 6, id)
	return fmt.Sprintf("code%s", content)
}
