package database

import (
	"api/services/util/log"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

var engines *xorm.EngineGroup

// ConnectionConfig 連線設定
type ConnectionConfig struct {
	Charset  		string
	Port	     	string
	MasterHostname 	string
	MasterUsername 	string
	MasterPassword 	string
	MasterDbName   	string
	SlaveHostname 	string
	SlaveUsername 	string
	SlavePassword 	string
	SlaveDbName   	string
}

type MysqlSession struct {
	Engine *xorm.Engine
	Session *xorm.Session
	Close func()
}

// Validate 驗證連線設定是否有錯
func (c *ConnectionConfig) Validate() error {
	return nil
}

// 連上Mysql
func connectMySQL(c *ConnectionConfig) *xorm.EngineGroup {
	var err error
	if c.Charset == "" {
		c.Charset = "utf8mb4"
	}
	connects := []string {
		c.MasterUsername+":"+c.MasterPassword+"@tcp("+c.MasterHostname+":"+c.Port+")/"+c.MasterDbName+"?charset="+c.Charset+"&parseTime=true",
		c.SlaveUsername+":"+c.SlavePassword+"@tcp("+c.SlaveHostname+":"+c.Port+")/"+c.SlaveDbName+"?charset="+c.Charset+"&parseTime=true",
	}
	engine, err := xorm.NewEngineGroup("mysql", connects)
	if err != nil {
		log.Error("database engine error", err)
		panic(err)
	}
	ENV := viper.GetString("ENV")
	if ENV != "prod" {
		//engine.ShowExecTime(true)
		engine.ShowSQL(true)
	}
	return engine
}

// GetEngine 取得某個 database 的 engine 連線
func GetMysqlEngineGroup() *xorm.EngineGroup {
	// get db config
	dbConfig := &ConnectionConfig{
		Port:     viper.GetString("database.mysql.Port"),
		Charset:  viper.GetString("database.mysql.Charset"),
		MasterHostname: viper.GetString("database.mysql.master.Hostname"),
		MasterUsername: viper.GetString("database.mysql.master.Username"),
		MasterPassword: viper.GetString("database.mysql.master.Password"),
		MasterDbName:   viper.GetString("database.mysql.master.dbName"),
		SlaveHostname: viper.GetString("database.mysql.slave.Hostname"),
		SlaveUsername: viper.GetString("database.mysql.slave.Username"),
		SlavePassword: viper.GetString("database.mysql.slave.Password"),
		SlaveDbName:   viper.GetString("database.mysql.slave.dbName"),
	}
	//log.Debug("database config", dbConfig)
	// connect to db
	engines := connectMySQL(dbConfig)
	return engines
}

func Mysql() *xorm.Engine {
	engine := GetMysqlEngineGroup()
	return engine.Master()
}

func GetMasterEngine() *xorm.Session {
	engine := GetMysqlEngineGroup()
	return engine.Master().NewSession()
}

func GetSlaveEngine() *xorm.Engine {
	engine := GetMysqlEngineGroup()
	return engine.Slave()
}


//取Mysql Engine
func GetMysqlEngine() *MysqlSession {
	engine := GetMysqlEngineGroup()
	mysqlEngine := &MysqlSession{
		Engine:  engine.Slave(),
		Session: engine.Master().NewSession(),
		Close: func() {
			engine.Master().Close()
			engine.Slave().Close()
		},
	}
	return mysqlEngine
}



