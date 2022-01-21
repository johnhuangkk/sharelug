package TransferService

import (
	"api/services/Enum"
	"api/services/dao/transfer"
	"api/services/database"
	"api/services/entity"
	"api/services/util/KgiAtmBank"
	"api/services/util/SCSBank"
	"api/services/util/log"
	"api/services/util/tools"
	"fmt"
	"github.com/spf13/viper"
)

//建立轉帳
func CreateTransfer(engine *database.MysqlSession, OrderId string, amount int64, Type string) (entity.TransferData, error) {
	var Entity entity.TransferData
	BankName, BankCode, account, err := createBankTransfer(int(amount),809, Type)
	log.Debug("sss", BankName, account)
	if err != nil {
		return Entity, err
	}
	Entity.OrderId = OrderId
	Entity.BankAccount = account
	Entity.BankName = BankName
	Entity.BankCode = BankCode
	Entity.Amount = amount
	Entity.Currency = "NT"
	Entity.TransType = Type
	Entity.ExpireDate = tools.GenerateTransferExpireTime(2)
	Entity.TransferStatus = Enum.TransferInit
	_, err = transfer.InsertTransfer(engine, Entity)
	if err != nil {
		log.Error("insert transfer data Error", err)
		return Entity, err
	}
	return Entity, nil
}

// Inner
func createBankTransfer(amount, bank int, transType string) (bankName, bankCode, account string, err error) {
	if bank == 11 {
		bankName = viper.GetString("ScsBank.BankName")
		bankCode = viper.GetString("ScsBank.BankCode")
		prefixCode := viper.GetString("ScsBank.VAPrefixCode")
		account,err = SCSBank.GenerateAtmVirtualAccount(prefixCode,amount,3)
		return bankName, bankCode, account, nil
	}
	if bank == 809 {
		bankName = viper.GetString("KgiBank.BankName")
		bankCode = viper.GetString("KgiBank.BankCode")
		account, err = KgiAtmBank.GenerateAtmVirtualAccount(amount, transType)
		return bankName, bankCode, account, nil
	}
	return "", "", "", fmt.Errorf("bank not Found")
}