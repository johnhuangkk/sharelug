package api

import (
	"api/config/middleware"
	"bytes"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

var configFile string

type Website struct {
	*gin.Engine
}

func init() {
	flag.StringVar(&configFile, "c", "config/config.yaml", "Configuration file path.")
	flag.Parse()
}

func NewDevelopment() *Website {
	s := &Website{
		Engine: gin.New(),
	}
	//讀取CONFIG
	err := readConfig()
	if err != nil {
		panic(err)
	}
	commonSetting(s)

	return s
}

//讀取config的檔案
func readConfig() error {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = viper.ReadConfig(bytes.NewBuffer(content))
	if err != nil {
		return err
	}

	return nil
}

func commonSetting(s *Website) {


	//var cookieName = viper.GetString("cookie.name")
	//使用redis
	//hostName := viper.GetString("redis.hostname")
	//store, err := redis.NewStore(10, "tcp", hostName, "", []byte("secret"))
	//if err != nil {
	//	panic(err)
	//}
	//使用cookie
	//store := cookie.NewStore([]byte("secret"))

	//s.Use(sessions.Sessions(cookieName, store))
	s.Use(middleware.Logger(), gin.Recovery())
	s.Use(middleware.Middleware(), gin.Recovery())
	//image
	s.Static("/static", "./www/static")
	//html
	s.LoadHTMLGlob("views/**/*")

	s.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusNotFound, "index/404.html", gin.H{"user": ""})
	})
}

func (s *Website) WebRun() {
	// set graceful shutdown
	zap.S().Fatal(s.Run(":8080"))
}
