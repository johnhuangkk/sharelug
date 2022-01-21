package controllers

import (
	"api/services/Service/KgiBank"
	"api/services/errorMessage"
	"api/services/model"
	"api/services/util/log"
	"api/services/util/response"
	"api/services/util/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"time"
)
//匯出KGI EACH提領
func ExporterKgiEachAction(ctx *gin.Context) {
	resp := response.New(ctx)
	filename, err := KgiBank.ExporterKgiEachFile()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}

	file, err := os.Open(filename) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename=EACH"+time.Now().Format("20060102")+".txt")
	ctx.Writer.Header().Add("Content-type", "application/txt")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		//resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}
//匯出KGI ACH提領
func ExporterKgiAchAction(ctx *gin.Context) {
	//resp := response.New(ctx)
	filename, err := KgiBank.ExporterKgiAchFile()
	if err != nil {
		//resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}

	file, err := os.Open(filename) //Create a file
	if err != nil {
		//resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename=ACH"+time.Now().Format("20060102")+".txt")
	ctx.Writer.Header().Add("Content-type", "application/txt")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		//resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}
//匯出使用者餘額列表
func GetMemberReportAction(ctx *gin.Context) {
	resp := response.New(ctx)
	filename, err := model.HandleMemberReportExporter()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	file, err := os.Open(filename) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+time.Now().Format("20060102")+".xlsx")
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}
//匯出訂單報表 是否有扣除服務費
func GetOrderReportAction(ctx *gin.Context) {
	resp := response.New(ctx)
	filename, err := model.HandleGetOrderReportExporter()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	file, err := os.Open(filename) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+time.Now().Format("20060102")+".xlsx")
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}

func ExporterDayStatementAction(ctx *gin.Context) {
	day := ctx.Param("day")
	log.Debug("day", day)
	path := tools.GetFilePath("/temp/day/", "", 0)
	filename := fmt.Sprintf("Day%s.%s", day, "xlsx")
	file, err := os.Open(fmt.Sprintf("%s%s", path, filename)) //Create a file
	if err != nil {
		log.Error("Open file Error", err)
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename=" + filename)
	ctx.Writer.Header().Add("Content-type", "application/txt")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		return
	}
}
//匯出使用者餘額列表
func ExporterInvoiceReportAction(ctx *gin.Context) {
	resp := response.New(ctx)
	day := ctx.Param("day")
	filename, err := model.HandleInvoiceReportExporter(day[:3], day[3:])
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	file, err := os.Open(filename) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+time.Now().Format("20060102")+".xlsx")
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}
//取出會員發票資料
func ExporterUserInvoiceReportAction(ctx *gin.Context) {
	resp := response.New(ctx)
	UserId := ctx.Param("id")
	filename, err := model.HandleUserInvoiceReportExporter(UserId)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	file, err := os.Open(filename) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+time.Now().Format("20060102")+".xlsx")
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}
//匯出銀行賣家資料
func ExporterBankReportAction(ctx *gin.Context) {
	resp := response.New(ctx)
	filename, err := model.HandleGetBankReportExporter()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	file, err := os.Open(filename) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+time.Now().Format("20060102")+".xlsx")
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}

func ExporterBalancesReportAction(ctx *gin.Context) {
	resp := response.New(ctx)
	filename, err := model.HandleGetBalancesReportExporter()
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	file, err := os.Open(filename) //Create a file
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
	defer file.Close()
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+time.Now().Format("20060102")+".xlsx")
	ctx.Writer.Header().Add("Content-type", "application/octet-stream")
	_, err = io.Copy(ctx.Writer, file)
	if err != nil {
		resp.Fail(errorMessage.GetMessageByCode(1001001)).Send()
		return
	}
}

