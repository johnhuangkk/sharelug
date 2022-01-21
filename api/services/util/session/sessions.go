package session

import (
	"api/services/util/log"
	"api/services/util/redis"
	"encoding/json"
	"time"
)

type Session struct {
	Name   string
	TTL    int64 // seconds
}

func NewUser(session string) *Session {
	return &Session{
		Name: "user" + session,
		TTL: 0,
	}
}

func OldUser(session string) *Session {
	return &Session{
		Name: "old" + session,
		TTL: 600,
	}
}

func Bind(session string) *Session {
	return &Session{
		Name: "bind" + session,
		TTL: 60000,
	}
}

func (sess *Session) Put(key string, value interface{}) error {
	var content string
	data := make(map[string]interface{})

	if redis.New().Exists(sess.Name) == true {
		content = redis.New().GetRedis(sess.Name)
	} else {
		content = "{}"
	}
	json.Unmarshal([]byte(content), &data)
	data[key] = value
	bytes, _ := json.Marshal(data)
	err := redis.New().SetRedis(sess.Name, string(bytes), time.Duration(sess.TTL) * time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (sess *Session) Get(key string) interface{} {
	var data map[string]interface{}
	var content string
	content = redis.New().GetRedis(sess.Name)
	json.Unmarshal([]byte(content), &data)
	return data[key]
}

func (sess *Session) Removes(key string) {
	var data map[string]interface{}
	var content string
	content = redis.New().GetRedis(sess.Name)
	json.Unmarshal([]byte(content), &data)
	delete(data, key)
	return
}

func (sess *Session) Destroy() error {
	err := redis.New().DelRedis(sess.Name)
	if err != nil {
		return err
	}
	return nil
}

//取session
func GetSession(uuid string, key string) string {
	sess := Session{
		Name:"user" + uuid,
		TTL: 86400,
	}
	uid := sess.Get(key)
	if uid != nil {
		return uid.(string)
	}
	return ""
}

//取session
func GetOldSession(uuid string, key string) string {
	sess := Session{
		Name:"old" + uuid,
		TTL: 0,
	}
	uid := sess.Get(key)
	if uid != nil {
		log.Debug("session get uid =>", uid)
		return uid.(string)
	}
	return ""
}
//
////設session
//func SetSession(uid string, ctx *gin.Context) error {
//	session := sessions.Default(ctx)
//	session.Set(userkey, uid)
//	err := session.Save()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
////刪除session
//func DeleteSession(ctx *gin.Context) error {
//	session := sessions.Default(ctx)
//	session.Delete(userkey)
//	err := session.Save()
//	if err != nil {
//		log.Error("session save Error!!", err)
//		return err
//	}
//	return nil
//}


