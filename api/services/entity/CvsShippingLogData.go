package entity

import (
	"api/services/Enum"
	"fmt"
)

type CvsShippingLogData struct {
	Id         int64  `xorm:"pk int(10) unique autoincr"`
	ShipNo     string `xorm:"varchar(20) comment('運送單號')"`
	CvsType    string `xorm:"varchar(20) comment('超商')"`
	Type       string `xorm:"varchar(10) comment('寫入類型 API/XML 若為XML 會寫入 R22 相關系列')"`
	Text       string `xorm:"varchar(50) comment('前端顯示文字')"`
	IsShow     bool   `xorm:"tinyint(1) NOT NULL comment('是否顯示到前端： 1 顯示 / 0 不顯示')"`
	DateTime   string `xorm:"datetime notnull comment('交易日期')"`
	CreateTime string `xorm:"datetime notnull comment('建立日期')"`
	Log        string `xorm:"text comment('原始資訊')"`
	FileName   string `xorm:"varchar(50) comment('檔案名稱')"`
}

func (cS *CvsShippingLogData) SetLogDataText(statusDetail string) {

	if cS.Type == `RS9` {
		cS.Text = setUnusualStatusText(statusDetail)
	} else {
		switch cS.CvsType {
		case Enum.CVS_7_ELEVEN:
			cS.Text = setSevenStatusText(cS.Type, statusDetail)
		case Enum.CVS_HI_LIFE:
			cS.Text = setHiLifStatusText(cS.Type, statusDetail)
		case Enum.CVS_FAMILY:
			cS.Text = setFamilyStatusText(cS.Type, statusDetail)
		case Enum.CVS_OK_MART:
			cS.Text = setOKStatusText(cS.Type, statusDetail)
		}
	}
}
func setSevenStatusText(code string, statusDetail string) string {
	mapStatusText := map[string]string{
		"000": "成功交寄",
		"011": "作業錯誤",
		"012": "車輛故障",
		"013": "天候不佳",
		"014": "道路中斷",
		"015": "門市停業中",
		"016": "缺件",
		"017": "門市報缺",
		"101": "門市配達",
		"102": "EC 管制品配達",
		"018": "寄件貨態異常協尋中",
		"019": "取件包裹異常協尋中",
		"201": "退貨門市配達",
		"202": "交貨便收件",
		"203": "退貨便收件",
		"204": "異常收退",
		"301": "取消寄件",
		"302": "寄件遺失進行賠償程序",
		"303": "取件遺失進行賠償程序",
		"777": "已完成取件",
		"888": "賣家完成取件",
		// CEIN
		"00":  "進驗成功",
		"09":  "未到貨",
		"31":  "商品破損",
		"32":  "超才",
		"33":  "違禁品",
		"34":  "訂單資料重複",
		"35":  "已過門市進貨日 (貨到物流中心的時間已過商品出貨日)",
		"36":  "門市關轉",
		"37":  "條碼規格錯誤",
		"38":  "條碼無法判讀",
		"39":  "條碼資料錯誤",
		"60":  "物流中心理貨中*",
		"61":  "商品遺失*",
		"62":  "門市不配送*",
		"63":  "包裹異常不配送*",
		"64":  "取消寄件再次寄送(直接轉 C 店) *",
		"65":  "提早轉 C 店-廠商因素(直接轉 C 店) *",
		"66":  "提早轉 C 店-超商因素(直接轉 C 店) *",
		"R00": "進驗成功",
		"R09": "未到貨",
		"R31": "商品破損",
		"R32": "超才",
		"R33": "違禁品",
		"R34": "訂單資料重複",
		"R35": "已過門市進貨日 (貨到物流中心的時間已過商品出貨日)",
		"R36": "門市關轉",
		"R37": "條碼規格錯誤",
		"R38": "條碼無法判讀",
		"R39": "條碼資料錯誤",
		"R60": "物流中心理貨中*",
		"R61": "商品遺失*",
		"R62": "門市不配送*",
		"R63": "包裹異常不配送*",
		"R64": "取消寄件再次寄送(直接轉 C 店) *",
		"R65": "提早轉 C 店-廠商因素(直接轉 C 店) *",
		"R66": "提早轉 C 店-超商因素(直接轉 C 店) *",
		//CERT
		"RT01": "未取退回物流中心",
		"RT11": "商品瑕疵",
		"RT12": "門市關店",
		"RT13": "門市轉店",
		"RT14": "廠商要求",
		"RT15": "違禁品",
		"RT21": "刷A給B",
		"RT22": "消費者要求",
		//CEDR
		"DR01": "未取退回物流中心",
		"DR11": "商品瑕疵",
		"DR12": "門市關店",
		"DR13": "門市轉店",
		"DR14": "廠商要求",
		"DR15": "違禁品",
		"DR21": "刷A給B",
		"DR22": "消費者要求",
	}
	return mapStatusText[code]
}
func setOKStatusText(typeX, statusDetail string) string {

	mappingCvsShippingText := map[string]string{
		"F27": "已完成交寄",
		"F84": "前往物流中心",
		"F71": "大物流驗收", "F63": "小物流驗收", "F03": "物流驗收",
		"F67": "物流驗退", "F07": "物流驗退",
		"F44": "配達取件店鋪", "F64": "配達取件店鋪", "F04": "配達取件店鋪",
		"F17": "已完成取件", "F65": "已完成取件", "F05": "已完成取件",
	}

	return mappingCvsShippingText[typeX]
}

func setHiLifStatusText(typeX, statusDetail string) string {
	fmt.Println("typeX: ", typeX)
	fmt.Println("statusDetail: ", statusDetail)

	mappingCvsShippingText := map[string]string{
		"R27": "已完成交寄", "R22": "已完成交寄",
		"R04": "前往物流中心",
		"R28": "配達取件店鋪", "RS4": "配達取件店鋪",
		"R29": "已完成取件",
		"R08": "未取退回物流中心",
		"RS9": "驗收異常",

		"Switch": "閉轉通知",
	}

	fmt.Println("mappingCvsShippingText", mappingCvsShippingText[typeX])

	return mappingCvsShippingText[typeX]
}

func setFamilyStatusText(typeX, statusDetail string) string {
	mappingCvsShippingText := map[string]string{
		"R22": "已完成交寄", "R22Api": "已完成交寄",
		"R23": "前往物流中心", "R23Api": "前往物流中心",
		"R25": "進行配送中",
		"R27": "物流出貨通知",
		"R08": "未取退回物流中心",

		"R04": "物流出貨",
		"RS9": "驗收異常",

		"RS4": "配達取件店鋪", "R28": "配達取件店鋪", "R28Api": "配達取件店鋪",
		"R96": "已完成取件", "R29": "已完成取件", "R29Api": "已完成取件",

		"Switch": "閉轉通知",
	}

	return mappingCvsShippingText[typeX]
}

// 設定特殊狀態
func setUnusualStatusText(status string) string {
	mappingCvsUnusualShippingText := map[string]string{
		"D04": "包裝廠包裝不良",
		"N05": "門市遺失",
		"S03": "物流遺失",
		"S06": "物流破損",
		"S07": "門市反應商品包裝不良",
		"T00": "正常驗退", // (只有ＯＫ有）
		"T01": "閉店、整修、無路線路 順",
		"T03": "條碼錯誤",
		"T04": "條碼重複",
		"T08": "超材",
		"XXX": "未到貨-取消",
		"T05": "貨物進店後發生異常提早退貨",
	}

	return mappingCvsUnusualShippingText[status]
}
