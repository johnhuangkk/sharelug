package middleware

import (
	"api/services/util/tools"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"
)

var hostname string
var ClientIP string

func Logger() gin.HandlerFunc {

	logFilePath := viper.GetString("logFilePath")
	logFileName := viper.GetString("logFileName")
	fileName := path.Join(logFilePath, logFileName + "_" + time.Now().Format("20060102") + ".log") //日誌檔案
	src, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)//寫入檔案
	if err != nil {
		fmt.Println("err", err)
	}
	logger := logrus.New()                       //例項化
	logger.Out = src                             //設定輸出
	logger.SetLevel(logrus.DebugLevel)           //設定日誌級別
	logger.SetFormatter(&logrus.TextFormatter{}) //設定日誌格式

	return func(ctx *gin.Context) {
		var bodyBytes []byte
		if ctx.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(ctx.Request.Body)
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		hostname = ctx.Request.Host
		ClientIP = ctx.ClientIP()
		startTime := time.Now()               // 開始時間
		ctx.Next()                            // 處理請求
		endTime := time.Now()                 // 結束時間
		latencyTime := endTime.Sub(startTime) // 執行時間
		reqMethod := ctx.Request.Method       // 請求方式
		reqUri := ctx.Request.RequestURI      // 請求路由
		reqPost, _ := tools.JsonEncode(ctx.Request.PostForm)
		reqBody := string(bodyBytes)
		statusCode := ctx.Writer.Status() // 狀態碼
		clientIP := GetClientIP()       // 請求IP
		var heading bytes.Buffer
		for k, v := range ctx.Request.Header {
			head := make(map[string]interface{})
			head[k] = v
			jsonString, _ := json.Marshal(head)
			heading.WriteString(string(jsonString))
		}
		inf, _ := net.Interfaces()
		// 日誌格式
		logger.Infof("| %3d | %13v | %15s | %s | %s | post=[%s] | body=[%s] | heading=[%s] | inf=[%v]", statusCode, latencyTime, clientIP, reqMethod, reqUri, reqPost, reqBody, heading.String(), inf)
	}
}

func GetHostname() string {
	return hostname
}

func GetClientIP() string {
	return ClientIP
}