package model

import (
	"api/services/Service/SCSBank"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"github.com/spf13/viper"
)

// 上海銀行入金檢查
// 流程: 向銀行端查詢金流資料 -> 本地寫入金流資料 -> 向銀行端更新金流資料索引 -> 更新訂單狀態
func FetchAccountingRecord() {
	log.Debug("FetchAccountingRecord[Begin]")
	host := viper.GetString("ScsBank.VAFetchHost")
	client := SCSBank.NewClient(host)

	rs, err := client.GetQueryAccountingData()
	if err != nil {
		log.Debug("FetchAccountingRecord1:",err.Error())
	}

	entities := privateGenerateRecordToEntity(rs)
	result := privateInsertEntities(entities)

	// TODO: 判斷寫入資料是否筆數正確，正確則更新銀行端索引
	if false && result {
		lastRecord := rs[len(rs)-1]
		privateUpdateBankData(client,lastRecord.RAccount, lastRecord.BTXDate, lastRecord.SeqNo)
	}

	for _,v := range entities {
		if err := HandleTransferForSCSBank(v); err != nil {
			log.Debug("FetchAccountingRecord2:",err.Error())
		}
	}
}

func privateUpdateBankData(client SCSBank.Client,racount, btxDate, seqNo string) bool {
	result, err := client.GetAckAccountingData(racount, btxDate, seqNo)
	if err != nil {
		log.Debug("FetchAccountingRecord3:",err.Error())
		return false
	}

	if !result {
		log.Debug("FetchAccountingRecord3:","Fail")
		return false
	}
	return true
}

func privateGenerateRecordToEntity(records []SCSBank.IncomeRecord) (entities []entity.ScsBankAccountData) {
	for _,r := range records {
		txKey := r.RAccount + r.BTXDate + r.SeqNo
		valid := 0
		if r.Valid {
			valid = 1
		}
		e := entity.ScsBankAccountData{
			TxKey:    txKey,
			RAccount: r.RAccount,
			BtxDate:  r.BTXDate,
			SeqNo:    r.SeqNo,
			SAccount: r.SAccount,
			BusDate:  r.BUSDate,
			TxDate:   r.TXDate,
			TxTime:   r.TXTime,
			VAccount: r.VAccount,
			Amount:   r.Amount,
			DC:       r.DC,
			Notes:    r.Notes,
			CheckSum: r.Checksum,
			IsValid:  valid,
			RawData:  r.Raw,
		}
		entities = append(entities,e)
	}
	return
}

func privateInsertEntities(entities []entity.ScsBankAccountData) bool {
	// 將紀錄寫入DB
	engine := database.GetMysqlEngine()
	defer engine.Close()
	insertCount ,err := engine.Session.Insert(entities)
	if err != nil {
		log.Debug("insertRecords:",err.Error())
	}
	return insertCount == int64(len(entities))
}
