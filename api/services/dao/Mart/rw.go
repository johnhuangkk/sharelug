package Mart

import (
	"api/services/Enum"
	"api/services/VO/FamilyMart"
	"api/services/database"
	"api/services/entity"
	"api/services/util/export"
	"api/services/util/log"
	"api/services/util/tools"
	"errors"
	"fmt"
	"time"
)

// Insert
func WriteInsertShippingOrder(ship interface{}) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()

	count, err := engine.Session.InsertOne(ship)
	if err != nil {
		return errors.New(fmt.Sprintf("insertShippingOrder Err:%v", err.Error()))
	}

	if count != 1 {
		return errors.New(fmt.Sprintf("insertShippingOrder NotOne:%d", count))
	}
	return nil
}

// Select
func ReadShippingOrder(ship *entity.MartOkShippingData) error {
	r, e := database.GetMysqlEngineGroup().Get(ship)
	if e != nil {
		return e
	}
	if !r {
		return errors.New("not found")
	}
	return nil
}

// Update
func WriteFamilyShipStateByEcOrderNo(ecOrderNo, state, stateCode, record, lg string) error {
	sql := "UPDATE mart_family_shipping_data SET record=CONCAT(?,record), log=CONCAT(?,log), state=?, state_code=? WHERE ec_order_no=?"
	return privateExecUpdateOne(sql, record, lg, state, stateCode, ecOrderNo)
}

func WriteFamilyLogByEcOrderNo(ecOrderNo, lg string) error {
	sql := "UPDATE mart_family_shipping_data SET log=CONCAT(?,log) WHERE ec_order_no=?"
	return privateExecUpdateOne(sql, lg, ecOrderNo)
}

func WriteFamilyShipStateByShipNo(shipNo, state, stateCode, record, lg string) error {
	sql := "UPDATE mart_family_shipping_data SET record=CONCAT(?,record),log=CONCAT(?,log),state=?, state_code=? WHERE ship_no=?"
	return privateExecUpdateOne(sql, record, lg, state, stateCode, shipNo)
}

func WriteFamilyShipStateWithNeedChangeByEcOrderNo(ecOrderNo, state, stateCode, record, lg string, needChange bool) error {
	sql := "UPDATE mart_family_shipping_data SET record=CONCAT(?,record), log=CONCAT(?,log), state=?, state_code=?, need_change=? WHERE ec_order_no=?"
	return privateExecUpdateOne(sql, record, lg, state, stateCode, needChange, ecOrderNo)
}

func WriteFamilyShipStateWithNeedChangeByShipNo(shipNo, state, stateCode, record, lg string, needChange bool) error {
	sql := "UPDATE mart_family_shipping_data SET record=CONCAT(?,record),log=CONCAT(?,log),state=?, state_code=?, need_change=? WHERE ship_no=?"
	return privateExecUpdateOne(sql, record, lg, state, stateCode, needChange, shipNo)
}

func WriteOKShipStateWithRecord(shipNo, state, stateCode, record string, lg string) error {
	sql := "UPDATE mart_ok_shipping_data SET record=CONCAT(?,record),log=CONCAT(?,log),state=?, state_code=? WHERE ship_no=?"
	return privateExecUpdateOne(sql, record, lg, state, stateCode, shipNo)
}

func WriteHiLifeShipStateByShipNo(shipNo, state, stateCode, record string, lg string) error {
	sql := "UPDATE mart_hi_life_shipping_data SET record=CONCAT(?,record),log=CONCAT(?,log),state=?, state_code=? WHERE ship_no=?"
	return privateExecUpdateOne(sql, record, lg, state, stateCode, shipNo)
}

func WriteOKShipStateWithNeedChangeByShipNo(shipNo, state, stateCode, record string, needChange bool, lg string) error {
	sql := "UPDATE mart_ok_shipping_data SET record=CONCAT(?,record),log=CONCAT(?,log),state=?, state_code=?, need_change=? WHERE ship_no=?"
	return privateExecUpdateOne(sql, record, lg, state, stateCode, needChange, shipNo)
}

func InsertFileFamily(parentId, eshopId, fType, fId, fDate, content string) error {
	return privateInsertFile(Enum.FAMILY, parentId, eshopId, fType, fId, fDate, content)
}

func InsertFileOK(parentId, eshopId, fType, fId, fDate, content string) error {
	return privateInsertFile(Enum.OK, parentId, eshopId, fType, fId, fDate, content)
}

func InsertFileHiLife(parentId, eshopId, fType, fId, fDate, content string) error {
	return privateInsertFile(Enum.HILIFE, parentId, eshopId, fType, fId, fDate, content)
}

// 更新 貨態更新時間
func UpdateOkToNow(shipNo string) error {
	condi := entity.MartOkShippingData{ShipNo: shipNo}
	data := entity.MartOkShippingData{UpdateDT: tools.Now("YmdHis")}
	return privateUpdateOne(&data, &condi)
}
func UpdateHiLifeToNow(shipNo string) error {
	condi := entity.MartHiLifeShippingData{ShipNo: shipNo}
	data := entity.MartHiLifeShippingData{UpdateDT: tools.Now("YmdHis")}
	return privateUpdateOne(&data, &condi)
}
func UpdateFamilyToNow(shipNo string) error {
	condi := entity.MartFamilyShippingData{ShipNo: shipNo}
	data := entity.MartFamilyShippingData{UpdateDT: tools.Now("YmdHis")}
	return privateUpdateOne(&data, &condi)
}

// 更新 閉轉有效時間
func UpdateOkSwitchTime(shipNo, switchDT string) error {
	condi := entity.MartOkShippingData{ShipNo: shipNo}
	data := entity.MartOkShippingData{SwitchDT: switchDT}
	return privateUpdateOne(&data, &condi)
}
func UpdateHiLifeSwitchTime(shipNo, switchDT string) error {
	condi := entity.MartHiLifeShippingData{ShipNo: shipNo}
	data := entity.MartHiLifeShippingData{SwitchDT: switchDT}
	return privateUpdateOne(&data, &condi)
}
func UpdateFamilySwitchTime(ecOrderNo, switchDT string) error {
	condi := entity.MartFamilyShippingData{EcOrderNo: ecOrderNo}
	data := entity.MartFamilyShippingData{SwitchDT: switchDT}
	return privateUpdateOne(&data, &condi)
}

// 更新訂單運送狀態
func UpdateOrderShippingStatus(ecOrderNo, status string) error {
	condi := entity.OrderData{OrderId: ecOrderNo}
	data := entity.OrderData{ShipStatus: status}
	return privateUpdateOne(&data, &condi)
}

// 更新運費
func UpdateFeeHiLife(shipNo, fee string) error {
	f := entity.MartHiLifeShippingData{ShipFee: fee}
	condition := entity.MartHiLifeShippingData{ShipNo: shipNo}
	return privateUpdateOne(&f, &condition)
}

func UpdateFeeFamilyByEcOrderNo(ecOrderNo, fee string) error {
	f := entity.MartFamilyShippingData{ShipFee: fee}
	condition := entity.MartFamilyShippingData{EcOrderNo: ecOrderNo}
	return privateUpdateOne(&f, &condition)
}

// 更新遺失貨態
func UpdateIsLoseHiLife(shipNo string, lose bool) error {
	f := entity.MartHiLifeShippingData{IsLose: lose}
	condition := entity.MartHiLifeShippingData{ShipNo: shipNo}
	return privateUpdateOne(&f, &condition)
}

func UpdateIsLoseFamilyByEcOrderNo(ecOrderNo string, lose bool) error {
	f := entity.MartFamilyShippingData{IsLose: lose}
	condition := entity.MartFamilyShippingData{EcOrderNo: ecOrderNo}
	return privateUpdateOne(&f, &condition)
}

// 更新為退貨狀態
func UpdateOnReturnFamily(ecOrderNo string, onReturn bool) error {
	f := entity.MartFamilyShippingData{OnReturn: onReturn}
	condition := entity.MartFamilyShippingData{EcOrderNo: ecOrderNo}
	return privateUpdateOne(&f, &condition)
}

func UpdateOnReturnHiLife(shipNo string, onReturn bool) error {
	f := entity.MartHiLifeShippingData{OnReturn: onReturn}
	condition := entity.MartHiLifeShippingData{ShipNo: shipNo}
	return privateUpdateOne(&f, &condition)
}

func UpdateOnReturnOK(shipNo string, onReturn bool) error {
	f := entity.MartOkShippingData{OnReturn: onReturn}
	condition := entity.MartOkShippingData{ShipNo: shipNo}
	return privateUpdateOne(&f, &condition)
}

// 更新 EshopId
func UpdateEshopIdOK(shipNo, ecNo string) error {
	f := entity.MartOkShippingData{EshopId: ecNo}
	condition := entity.MartOkShippingData{ShipNo: shipNo}
	return privateUpdateOne(&f, &condition)
}

func privateInsertFile(vendor, parentId, eshopId, fType, fId, fDate, content string) error {
	entity := entity.MartShippingFileData{
		Vendor:   vendor,
		ParentId: parentId,
		EshopId:  eshopId,
		FileType: fType,
		FileId:   fId,
		FileDate: fDate,
		Content:  content,
	}
	engine := database.GetMysqlEngine()
	defer engine.Close()
	c, err := engine.Session.InsertOne(entity)
	if err != nil {
		return err
	}

	if c != 1 {
		return errors.New("InsertFile count not fit")
	}

	return nil
}

/**
取得超商店家資訊
*/
func QueryCVSData(engine *database.MysqlSession, storeId string, ent interface{}) error {
	var err error

	switch fmt.Sprintf("%T", ent) {
	case "*entity.SevenMyshipShopData":
		_, err = engine.Engine.Table(entity.SevenMyshipShopData{}).Select("*").Where("store_id = ?", storeId).Get(ent)
	case "*entity.MartFamilyStoreData":
		_, err = engine.Engine.Table(entity.MartFamilyStoreData{}).Select("*").Where("store_id = ?", storeId).Get(ent)
	case "*entity.MartHiLifeStoreData":
		_, err = engine.Engine.Table(entity.MartHiLifeStoreData{}).Select("*").Where("store_id = ?", storeId).Get(ent)
	case "*entity.MartOkStoreData":
		_, err = engine.Engine.Table(entity.MartOkStoreData{}).Select("*").Where("store_id = ?", storeId).Get(ent)
	}
	if err != nil {
		fmt.Errorf("%s", "系統錯誤")
		return err
	}

	return nil
}

// 寫出 萊爾富JSON給前端頁面用
func WriteHiLifeStoreData(stores []entity.MartHiLifeStoreData) error {
	js := map[string]map[string][]entity.MartHiLifeStoreData{
		"台北市": map[string][]entity.MartHiLifeStoreData{},
		"新北市": map[string][]entity.MartHiLifeStoreData{},
		"基隆市": map[string][]entity.MartHiLifeStoreData{},
		"桃園市": map[string][]entity.MartHiLifeStoreData{},
		"新竹市": map[string][]entity.MartHiLifeStoreData{},
		"新竹縣": map[string][]entity.MartHiLifeStoreData{},
		"苗栗縣": map[string][]entity.MartHiLifeStoreData{},
		"台中市": map[string][]entity.MartHiLifeStoreData{},
		"彰化縣": map[string][]entity.MartHiLifeStoreData{},
		"南投縣": map[string][]entity.MartHiLifeStoreData{},
		"雲林縣": map[string][]entity.MartHiLifeStoreData{},
		"嘉義市": map[string][]entity.MartHiLifeStoreData{},
		"嘉義縣": map[string][]entity.MartHiLifeStoreData{},
		"台南市": map[string][]entity.MartHiLifeStoreData{},
		"高雄市": map[string][]entity.MartHiLifeStoreData{},
		"屏東縣": map[string][]entity.MartHiLifeStoreData{},
		"宜蘭縣": map[string][]entity.MartHiLifeStoreData{},
		"花蓮縣": map[string][]entity.MartHiLifeStoreData{},
		"台東縣": map[string][]entity.MartHiLifeStoreData{},
		"澎湖縣": map[string][]entity.MartHiLifeStoreData{},
		"金門縣": map[string][]entity.MartHiLifeStoreData{},
		"連江縣": map[string][]entity.MartHiLifeStoreData{},
	}
	for _, s := range stores {
		j, ok := js[s.City]
		if !ok {
			j = map[string][]entity.MartHiLifeStoreData{}
		}
		d, ok2 := j[s.District]
		if !ok2 {
			d = []entity.MartHiLifeStoreData{}
		}
		d = append(d, s)
		j[s.District] = d
		js[s.City] = j
	}

	cities := []FamilyMart.City{}

	for k, v := range js {
		discs := []FamilyMart.District{}
		for j, v2 := range v {
			addresss := []FamilyMart.Address{}

			for _, v3 := range v2 {
				addr := FamilyMart.Address{
					Name:      v3.StoreAddress,
					StoreId:   v3.StoreId,
					StoreName: v3.StoreName,
				}
				addresss = append(addresss, addr)
			}
			disc := FamilyMart.District{Name: j, Address: addresss}
			discs = append(discs, disc)
		}
		city := FamilyMart.City{Country: k, Districts: discs}
		cities = append(cities, city)
	}

	fileName := "HiLife_" + time.Now().Format("20060102") + ".json"
	return export.JsonToData2(cities, fileName)
}

// 寫出 全家店鋪JSON給前端頁面用
func WriteFamilyStoreData(stores []entity.MartFamilyStoreData) error {
	js := map[string]map[string][]entity.MartFamilyStoreData{
		"台北市": map[string][]entity.MartFamilyStoreData{},
		"新北市": map[string][]entity.MartFamilyStoreData{},
		"基隆市": map[string][]entity.MartFamilyStoreData{},
		"桃園市": map[string][]entity.MartFamilyStoreData{},
		"新竹市": map[string][]entity.MartFamilyStoreData{},
		"新竹縣": map[string][]entity.MartFamilyStoreData{},
		"苗栗縣": map[string][]entity.MartFamilyStoreData{},
		"台中市": map[string][]entity.MartFamilyStoreData{},
		"彰化縣": map[string][]entity.MartFamilyStoreData{},
		"南投縣": map[string][]entity.MartFamilyStoreData{},
		"雲林縣": map[string][]entity.MartFamilyStoreData{},
		"嘉義市": map[string][]entity.MartFamilyStoreData{},
		"嘉義縣": map[string][]entity.MartFamilyStoreData{},
		"台南市": map[string][]entity.MartFamilyStoreData{},
		"高雄市": map[string][]entity.MartFamilyStoreData{},
		"屏東縣": map[string][]entity.MartFamilyStoreData{},
		"宜蘭縣": map[string][]entity.MartFamilyStoreData{},
		"花蓮縣": map[string][]entity.MartFamilyStoreData{},
		"台東縣": map[string][]entity.MartFamilyStoreData{},
		"澎湖縣": map[string][]entity.MartFamilyStoreData{},
		"金門縣": map[string][]entity.MartFamilyStoreData{},
		"連江縣": map[string][]entity.MartFamilyStoreData{},
	}

	for _, s := range stores {
		j, ok := js[s.City]
		if !ok {
			j = map[string][]entity.MartFamilyStoreData{}
		}
		d, ok2 := j[s.District]
		if !ok2 {
			d = []entity.MartFamilyStoreData{}
		}
		d = append(d, s)
		j[s.District] = d
		js[s.City] = j
	}
	var cities []FamilyMart.City

	for k, v := range js {
		var discs []FamilyMart.District
		for j, v2 := range v {
			var addresss []FamilyMart.Address

			for _, v3 := range v2 {
				addr := FamilyMart.Address{
					Name:      v3.StoreAddress,
					StoreId:   v3.StoreId,
					StoreName: v3.StoreName,
				}
				addresss = append(addresss, addr)
			}
			disc := FamilyMart.District{Name: j, Address: addresss}
			discs = append(discs, disc)
		}
		city := FamilyMart.City{Country: k, Districts: discs}
		cities = append(cities, city)
	}

	fileName := "FamilyMart_" + time.Now().Format("20060102") + ".json"
	return export.JsonToData2(cities, fileName)
}

// 寫出 OK店鋪JSON給前端頁面用
func WriteOkStoreData(stores []entity.MartOkStoreData) error {
	js := map[string]map[string][]entity.MartOkStoreData{
		"台北市": map[string][]entity.MartOkStoreData{},
		"新北市": map[string][]entity.MartOkStoreData{},
		"基隆市": map[string][]entity.MartOkStoreData{},
		"桃園市": map[string][]entity.MartOkStoreData{},
		"新竹市": map[string][]entity.MartOkStoreData{},
		"新竹縣": map[string][]entity.MartOkStoreData{},
		"苗栗縣": map[string][]entity.MartOkStoreData{},
		"台中市": map[string][]entity.MartOkStoreData{},
		"彰化縣": map[string][]entity.MartOkStoreData{},
		"南投縣": map[string][]entity.MartOkStoreData{},
		"雲林縣": map[string][]entity.MartOkStoreData{},
		"嘉義市": map[string][]entity.MartOkStoreData{},
		"嘉義縣": map[string][]entity.MartOkStoreData{},
		"台南市": map[string][]entity.MartOkStoreData{},
		"高雄市": map[string][]entity.MartOkStoreData{},
		"屏東縣": map[string][]entity.MartOkStoreData{},
		"宜蘭縣": map[string][]entity.MartOkStoreData{},
		"花蓮縣": map[string][]entity.MartOkStoreData{},
		"台東縣": map[string][]entity.MartOkStoreData{},
		"澎湖縣": map[string][]entity.MartOkStoreData{},
		"金門縣": map[string][]entity.MartOkStoreData{},
		"連江縣": map[string][]entity.MartOkStoreData{},
	}

	//fmt.Println(stores)
	for _, s := range stores {
		j, ok := js[s.City]
		if !ok {
			j = map[string][]entity.MartOkStoreData{}
		}
		d, ok2 := j[s.District]
		if !ok2 {
			d = []entity.MartOkStoreData{}
		}
		d = append(d, s)
		j[s.District] = d
		js[s.City] = j
	}
	cities := []FamilyMart.City{}
	//fmt.Println(js)
	for k, v := range js {
		discs := []FamilyMart.District{}
		for j, v2 := range v {
			addresss := []FamilyMart.Address{}

			for _, v3 := range v2 {
				addr := FamilyMart.Address{
					Name:      v3.StoreAddress,
					StoreId:   v3.StoreId,
					StoreName: v3.StoreName,
				}
				addresss = append(addresss, addr)
			}
			disc := FamilyMart.District{Name: j, Address: addresss}
			discs = append(discs, disc)
		}
		city := FamilyMart.City{Country: k, Districts: discs}
		cities = append(cities, city)
	}
	fmt.Println(cities)
	fileName := "OK_" + time.Now().Format("20060102") + ".json"
	return export.JsonToData2(cities, fileName)
}

func QueryHiLifeShippingData(ecOrderNo string) (data entity.MartHiLifeShippingData, err error) {
	data.EcOrderNo = ecOrderNo
	return data, queryData(&data)
}

func QueryOKShippingData(ecOrderNo string) (data entity.MartOkShippingData, err error) {
	data.EcOrderNo = ecOrderNo
	return data, queryData(&data)
}

func QueryOKShippingDataByShipNo(shipNo string) (data entity.MartOkShippingData, err error) {
	data.ShipNo = shipNo
	return data, queryData(&data)
}

func QueryFamilyShippingData(ecOrderNo string) (data entity.MartFamilyShippingData, err error) {
	data.EcOrderNo = ecOrderNo
	return data, queryData(&data)
}

func queryData(data interface{}) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	result, err := engine.Session.Get(data)
	if !result {
		return err
	}
	return nil
}

// Private
func privateExecUpdateOne(sqlOrArgs ...interface{}) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	result, err := engine.Session.Exec(sqlOrArgs...)
	if err != nil {
		log.Debug("privateExecUpdateOne:", err.Error())
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Debug("privateExecUpdateOne:", count, err.Error())
		return err
	}
	if count != 1 {
		log.Debug("privateExecUpdateOne:", count)
		return err
	}
	return nil
}

func privateUpdateOne(bean interface{}, condition interface{}) error {
	engine := database.GetMysqlEngine()
	defer engine.Close()
	c, err := engine.Session.Update(bean, condition)
	if err != nil {
		return err
	}
	if c != 1 {
		return errors.New("Update count not fit")
	}
	return nil
}
