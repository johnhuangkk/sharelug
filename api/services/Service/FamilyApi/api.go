package FamilyApi

import (
	"api/services/Enum"
	fm "api/services/Service/FamilyMartLogistics"
	"api/services/VO/FamilyMart"
	"api/services/dao/Cvs"
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 全家取號
func C2cOrderAdd(engine *database.MysqlSession, orderData entity.OrderData, sellerData entity.MemberData) (string, error) {
	var order FamilyMart.OrderAddRequest
	order.SetParams(orderData, sellerData)
	order.Products.Objs = []FamilyMart.OrderAddRequestProduct{{}}


	log.Info("order [%s]", order)

	var client fm.Client
	client.GetClient(false)

	resp, err := client.OrderAdd(order)
	if err != nil {
		log.Error("Family C2cOrderAdd Error: [%s]", err.Error())
		return ``, err
	}

	if resp.ErrorCode != "000" {
		return resp.OrderNo, fmt.Errorf(resp.ErrorMessage)
	}


	// 建立托運資訊
	var data entity.CvsShippingData
	data.InitInsert(Enum.CVS_FAMILY)

	data.ParentId = order.ParentId
	data.EcOrderNo = order.EcOrderNo
	data.ShipNo = resp.OrderNo
	data.ServiceType = order.GetServiceType()
	data.SenderName = order.SenderName
	data.SenderPhone = order.SenderPhone
	data.OriReceiverAddress = orderData.ReceiverAddress

	err = Cvs.InsertCvsShippingData(engine, data)

	if err != nil {
		log.Error("C2cOrderAdd InsertCvsShippingData data Error: [%]", data)
		log.Error("C2cOrderAdd InsertCvsShippingData Error: [%]", err.Error())
		return "", err
	}

	return data.ShipNo, nil
}

func PrintShippingOrder(orderNos []string) (imageData []byte, err error) {
	return fetchShippingOrderPrint("901", "0001", orderNos...)
}

func fetchShippingOrderPrint(parentId, eshopId string, orderNo ...string) (data []byte, err error) {
	req := fm.OrdersPrintRequest{
		ParentId: parentId, EshopId: eshopId,
		Orders: []fm.OrdersPrintRequestOrder{},
	}
	for _, v := range orderNo {
		req.Orders = append(req.Orders, fm.OrdersPrintRequestOrder{OrderNo: v})
	}

	var client fm.Client
	client.GetClient(false)

	path, r := client.OrderPrint(req)
	if !r {
		return nil, fmt.Errorf("%s", "Get print fail")
	}

	httpClient := http.Client{}
	resp, err := httpClient.Get(path)
	if err != nil {
		log.Error("fetchShippingOrderPrint:", err)
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("fetchShippingOrderPrint:", err)
		return nil, err
	}

	return data, nil
}

// 修改閉轉店
func SwitchStore(cvsShippingData entity.CvsShippingData, newStoreId string) (err error)  {

	var client fm.Client
	var resp fm.OrderSwitchResponse

	storeType := `2`
	if cvsShippingData.FlowType == `R` {
		storeType = `1`
	}

	client.GetClient(false)

	resp, err = client.OrderSwitch("901", "0001", cvsShippingData.ShipNo, cvsShippingData.EcOrderNo, newStoreId, storeType)
	if err != nil {
		return err
	}

	if len(resp.Content) == 0 {
		return fmt.Errorf("%s", "轉店異常")
	}

	if 	resp.Content[0].ErrorCode != "000" {
		return fmt.Errorf("%s", resp.Content[0].ErrorMessage)
	}

	return nil
}

func MartFamilySwitchStore(shipNo, ecOrderNo, newStoreId string, isReceiveStore bool) error {
	storeType := "1"
	if isReceiveStore {
		storeType = "2"
	}

	var client fm.Client
	client.GetClient(false)

	resp, err := client.OrderSwitch("901", "0001", shipNo, ecOrderNo, newStoreId, storeType)
	if err != nil {
		return err
	}

	if len(resp.Content) == 0 {
		return fmt.Errorf("%s", "轉店異常")
	}

	d := resp.Content[0]
	if d.ErrorCode != "000" {
		//state := "關轉失敗"
		//dirc := ""
		//detail := generateRecord("", "", dirc, state)
		//log := generateLog("", "", state, resp)
		//_ = Mart.WriteFamilyShipStateByShipNo(d.OrderNo, state, "1", detail, log)
		//return errors.New(state)
	}

	// 更新關轉等待時間為7天
	//Mart.UpdateFamilySwitchTime(ecOrderNo, "")
	//state := "關轉成功"
	//dirc := ""
	//detail := generateRecord("", "", dirc, state)
	//log := generateLog("", "", state, resp)
	//Mart.UpdateOrderShippingStatus(ecOrderNo, Enum.OrderShipOnShipping)
	//_ = Mart.WriteFamilyShipStateWithNeedChangeByShipNo(d.OrderNo, state, "1", detail, log, false)
	return nil
}