package redis

import (
	"api/services/util/log"
	"github.com/spf13/viper"
	"gopkg.in/redis.v5"
	"time"
)


type Redis struct {
	conn redis.Client
}

func New() *Redis {
	hostname := viper.GetString("redis.hostname")
	password := viper.GetString("redis.password")
	client := redis.NewClient(&redis.Options{
		Addr:     hostname,
		Password: password,
		DB:       0,
	})
	// 通過 cient.Ping() 來檢查是否成功連線到了 redis 伺服器
	_, err := client.Ping().Result()
	if err != nil {
		log.Error("Connection Redis error", err)
	}
	return &Redis{
		conn: *client,
	}
}

//存值
func (rc *Redis)SetRedis(key string, value string, exp time.Duration) error {
	defer rc.conn.Close()
	err := rc.conn.Set(key, value, exp).Err()
	log.Debug("redis set data", key, value, err)
	if err != nil {
		log.Error("redis set error", key, err)
		return err
	}
	return nil
}

//List存值
func (rc *Redis) SetListRedis(key string, value []string) error {
	defer rc.conn.Close()
	//選擇從左邊或是右邊 push 值進去 LPUSH 與 RPUSH
	for _, v := range value {
		err := rc.conn.RPush(key, v).Err()
		if err != nil {
			log.Error("redis set list error", err)
			return err
		}
	}
	return nil
}

func (rc *Redis) SetHashRedis(key string, field string, value string) error {
	defer  rc.conn.Close()
	err := rc.conn.HSet(key, field, value).Err()
	if err != nil {
		log.Error("redis set hash error", err)
		return err
	}
	return nil
}

func (rc *Redis) SetExpireRedis(key string, expire time.Duration) error {
	defer rc.conn.Close()
	err := rc.conn.Expire(key, expire).Err()
	if err != nil {
		log.Error("redis set expire error", err)
		return err
	}
	return nil
}

func (rc *Redis) GetHashRedis(key string, field string) (string, error) {
	defer rc.conn.Close()
	val, err := rc.conn.HGet(key, field).Result()
	if err != nil {
		//log.Error("redis get hash error", err)
		return "", err
	}
	return val, nil
}

func (rc *Redis) DelRedis(key string) error {
	defer rc.conn.Close()
	err := rc.conn.Del(key).Err()
	if err != nil {
		log.Error("redis delete hash key error", err)
		return err
	}
	return nil
}

func (rc *Redis) DelHashRedis(key string, field string) error {
	defer rc.conn.Close()
	err := rc.conn.HDel(key, field).Err()
	if err != nil {
		log.Error("redis delete hash key error", err)
		return err
	}
	return nil
}

func (rc *Redis) GetHashAllRedis(key string) (map[string]string, error) {
	defer rc.conn.Close()
	val, err := rc.conn.HGetAll(key).Result()
	if err != nil {
		log.Error("redis get hash all error", err)
		return nil, err
	}
	return val, nil
}

//取出指定的值
func (rc *Redis) GetRedis(key string) string {
	defer rc.conn.Close()
	val, err := rc.conn.Get(key).Result()
	if err != nil {
		log.Debug("get redis value error", key, err)
		return "{}"
	}
	return val
}

//取出指定範圍
func (rc *Redis) GetListRedis(key string) ([]string, error) {
	defer rc.conn.Close()
	//LRANGE 可以印出指定範圍的值，支援-1這種形式，表示最後一個值
	val, err := rc.conn.LRange(key, 0, -1).Result()
	if err != nil {
		log.Debug("get redis list value error", err)
		return nil, err
	}
	return val, nil
}

func (rc *Redis) Exists(key string) bool {
	defer rc.conn.Close()
	val := rc.conn.Exists(key).Val()
	return val
}

func (rc *Redis) Expire(key string, expire time.Duration) error {
	defer rc.conn.Close()
	val := rc.conn.Expire(key, expire)
	return val.Err()
}