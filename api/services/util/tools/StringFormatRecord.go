package tools

import (
	"api/services/VO/InvoiceVo"
	"api/services/VO/KgiBank"
	"api/services/util/log"
	"fmt"
	"reflect"
	"strings"
)

type FieldFormat struct {
	FieldName string
	FieldType string
	Length    int
	Pattern    string
}

type FieldFormats struct {
	fieldFormats []FieldFormat
}

func SetRule(fieldName string, fieldType string, length int, pattern string) *FieldFormat {
	fieldFormat := new(FieldFormat)
	fieldFormat.FieldName = fieldName
	fieldFormat.FieldType = fieldType
	fieldFormat.Length = length
	fieldFormat.Pattern = pattern
	return fieldFormat
}

func SetRules(fieldName string, fieldType string, length int, pattern string, fieldFormats *FieldFormats) {
	fieldFormat := SetRule(fieldName, fieldType, length, pattern)
	fieldFormats.fieldFormats = append(fieldFormats.fieldFormats, *fieldFormat)
}

func (fieldFormats *FieldFormats) ToString(body interface{}) string {

	var result []string

	for _, f := range fieldFormats.fieldFormats {
		key := reflect.ValueOf(body)
		value := key.FieldByName(f.FieldName)

		if "N" == f.FieldType {
			//N 主要為數字0~9 組成，排列方式為「右靠左補 0」
			content := fmt.Sprintf("%0*s", f.Length, value)
			result = append(result, content)
		} else if "S" == f.FieldType {
			//A 及 AN 主要為一般字元組成，排列方式為「左靠右空白」
			content := fmt.Sprintf("%-*s", f.Length, value)
			result = append(result, content)
		} else if "SC" == f.FieldType {
			content := AddSpace(f.Length, value.String())
			result = append(result, content)
		}
	}
	result = append(result, "\x0a")
	return strings.Join(result, "")
}

func AddSpace(length int, str string) string {
	charset := Utf8ToBig5(str)
	log.Debug("EncodeBig5", charset)
	result := ""
	for i:= 0;i < length - len(charset); i++ {
		result += " "
	}
	return charset + result
}

/**
 * read file
 * @param string $raw
 * @param bool $ignoreError
 */
func (fieldFormats *FieldFormats) SetRawData(raw string, body *KgiBank.Body) {
	start := 0
	for _, f := range fieldFormats.fieldFormats {
		length := start + f.Length
		if len(raw) > length {
			key := reflect.ValueOf(body)
			value := key.Elem().FieldByName(f.FieldName)
			value.SetString(Trim(raw[start:length]))
		}
		start += f.Length
	}
}

/**
 * read file
 * @param string $raw
 * @param bool $ignoreError
 */
func (fieldFormats *FieldFormats) SetInvoiceData(raw string, body *InvoiceVo.Awarded) {
	start := 0
	for _, f := range fieldFormats.fieldFormats {
		length := start + f.Length
		if len(raw) > length {
			key := reflect.ValueOf(body)
			value := key.Elem().FieldByName(f.FieldName)
			value.SetString(Trim(raw[start:length]))
		}
		start += f.Length
	}
}