package tools

import (
	"api/services/util/log"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"strings"
)

//Utf8 轉 Big5
func Utf8ToBig5(str string) string {
	big5, _, err := transform.String(traditionalchinese.Big5.NewEncoder(), str)
	if err != nil {
		log.Error("Utf8ToBig5: ", err)
	}
	return big5
}

//big5 轉 utf8
func Big5ToUtf8(str string) string {
	utf8, _, err := transform.String(traditionalchinese.Big5.NewDecoder(), str)
	if err != nil {
		log.Error("Utf8ToBig5: ", err)
	}
	return utf8
}

/**
Big5 轉 Utf8 by byte
*/
func Big5ToUtf8ByByte(data []byte) ([]byte, error) {
	big5ToUTF8 := traditionalchinese.Big5.NewDecoder()
	utf8, _, err := transform.Bytes(big5ToUTF8, data)
	if err != nil {
		log.Error("Utf8ToBig5: ", err)
	}
	return utf8, err
}

func Utf8ToUnicode(str string) string {
	uni := fmt.Sprintf("%U", []rune(str))
	step1 := strings.ReplaceAll(uni, "U+", "\\u")
	step2 := strings.ReplaceAll(step1, " ", "")
	step3 := strings.ReplaceAll(step2, "[", "")
	step4 := strings.ReplaceAll(step3, "]", "")

	return step4
}

func StructToJsonGetValueString(v interface{}) string {
	jsonData, _ := json.Marshal(v)
	jb := []byte(string(jsonData))
	var settings map[string]interface{}
	//var jsonValue string
	if err := json.NewDecoder(bytes.NewReader(jb)).Decode(&settings); err != nil {
		panic(err)
	}
	//for k, v := range settings {
	//
	//}

	return "123"
}
