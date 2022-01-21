package response

import (
	"api/services/util/log"
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	gCtx       *gin.Context
	message    string
	status     string
	statusCode int64
	errorName  string
	data       interface{}
}

type ErrorInterface interface {
	GetErrorName() string
	GetStatusCode() int64
	GetMessage() string
}

const (
	StatusSuccess  = "Success"
	StatusFail     = "Fail"
	StatusPanic    = "Panic"
	StatusError    = "Error"
	StatusConflict = "Conflict"
)

// New ...
func New(c *gin.Context) *Response {
	return &Response{
		gCtx: c,
	}
}

// Success ..
func (rsp *Response) Success(message string) *Response {
	rsp.statusCode = http.StatusOK
	rsp.status = StatusSuccess
	rsp.message = message

	return rsp
}

// Fail ...
func (rsp *Response) Fail(code int64, message string) *Response {

	log.Error("Fail Error,    " + message)
	rsp.statusCode = code
	rsp.status = StatusFail
	rsp.message = strings.Trim(message, "\b")

	return rsp
}

// Error ...
func (rsp *Response) Error(err interface{}) *Response {
	rsp.status = StatusError
	rsp.setErrorData(err)

	return rsp
}

// Panic ...
func (rsp *Response) Panic(err interface{}) *Response {
	rsp.status = StatusPanic
	rsp.setErrorData(err)

	return rsp
}

// Conflict ...
func (rsp *Response) Conflict(message string) *Response {
	rsp.status = StatusConflict
	rsp.statusCode = 409
	rsp.message = message

	return rsp
}

func (rsp *Response) setErrorData(err interface{}) {
	rsp.message = err.(error).Error()

	switch err.(type) {
	case ErrorInterface:
		error := err.(ErrorInterface)
		rsp.statusCode = error.GetStatusCode()
		rsp.errorName = error.GetErrorName()
	default:
		// default errors
		rsp.statusCode = 500
	}
}

// SetData ...
func (rsp *Response) SetData(d interface{}) *Response {
	rsp.data = d
	return rsp
}

func (rsp *Response) SendString() {

	sCode := http.StatusOK
	if rsp.status == StatusPanic {
		sCode = http.StatusInternalServerError
	}

	if rsp.status == StatusConflict {
		sCode = http.StatusConflict
	}

	resp := gin.H{
		"StatusCode": rsp.statusCode,
		"Status":     rsp.status,
		"Data":       rsp.data,
		"Message":    rsp.message,
	}

	rsp.gCtx.JSON(sCode, resp)
}

// Send ...
func (rsp *Response) Send() {
	sCode := http.StatusOK
	if rsp.status == StatusPanic {
		sCode = http.StatusInternalServerError
	}

	if rsp.status == StatusConflict {
		sCode = http.StatusConflict
	}

	resp := gin.H{
		"StatusCode": rsp.statusCode,
		"Status":     rsp.status,
		"Data":       rsp.data,
		"Message":    rsp.message,
	}

	if sCode != http.StatusOK {
		log.Debug("Response Data", resp, rsp.gCtx.Request.RequestURI, rsp.gCtx.ClientIP())
		rsp.gCtx.AbortWithStatusJSON(sCode, resp)
	} else {
		log.Debug("Response Data", resp, rsp.gCtx.Request.RequestURI, rsp.gCtx.ClientIP())
		rsp.gCtx.JSON(sCode, resp)
	}
}

// 輸出 xml
func (rsp *Response) XML(v interface{}) {
	body, _ := xml.Marshal(v)
	xml := xml.Header + string(body)
	rsp.gCtx.Data(http.StatusOK, `application/xml`, []byte(xml))
}
