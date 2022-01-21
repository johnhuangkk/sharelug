package Erp

import (
	"api/services/database"
	"api/services/entity"
	"api/services/util/log"
	"fmt"

	"api/services/util/redis"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const salt1 = `$h@re1ug`

func CreateErpUser(ctx *gin.Context) {
	engine := database.GetMysqlEngine().Session
	defer engine.Close()

	testpwd := `112345`

	password := []byte(testpwd + salt1)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return
	}
	uid, err := uuid.NewUUID()
	if err != nil {
		return
	}

	_, err = engine.Table(entity.ErpUser{}).Insert(entity.ErpUser{
		Uid:      uid.String(),
		Password: string(hashedPassword),
		Email:    `duke@sharelug.com`,
		Name:     `Duke`,
		Enable:   true,
	})
	if err != nil {
		return
	}
}

func Login(ctx *gin.Context) {
	var f interface{}
	ctx.ShouldBindJSON(&f)
	fmt.Println(f)
	rst, _ := ctx.GetRawData()
	fmt.Println(rst)
	email := ctx.Request.FormValue("email")
	pwd := ctx.Request.FormValue("password")
	flag, user := userExist(email)
	// fmt.Println(user)
	if flag {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd+salt1))
		if err == nil {
			sID, _ := uuid.NewUUID()
			//set cookie
			ctx.SetCookie(`erp_session`, sID.String(), 7200, "/", "", false, true)
			err := redis.New().SetRedis(sID.String(), user.Uid, 10800)
			if err != nil {
				return
			}
			//set session

		} else {
			fmt.Println(user.Password)
			fmt.Println(pwd + salt1)
			fmt.Println(err.Error())
			fmt.Println("Fail")
		}
	}

}

func userExist(email string) (bool, entity.ErpUser) {
	engine := database.GetMysqlEngine().Engine
	defer engine.Close()
	var data entity.ErpUser
	_, err := engine.Table(entity.ErpUser{}).Where(`email=?`, email).Get(&data)
	if err != nil {
		log.Error(err.Error())
		return false, data
	}
	return true, data

}

// func Logout(ctx *gin.Context) {
// 	email := ctx.Request.FormValue("email")
// 	pwd := ctx.Request.FormValue("password")
// }
