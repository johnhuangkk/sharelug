package tools

import (
	"api/services/util/log"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"time"
)

type Claims struct {
	UUid    string    `json:"uuid"`
	UserId  string    `json:"user_id"`
	StoreId string    `json:"store_id"`
	Ip      string    `json:"ip"`
	Exp     time.Time `json:"exp"`
}

type Condition struct {
	PersonId      string `json:"personId"`
	IdMark        string `json:"applyCode"`
	IdMarkDate    string `json:"applyYyymmdd"`
	IssueAreaCode string `json:"issueSiteId"`
}

var secret = []byte("SecretCode")

func GeneratorJWT(UUID string, UserId string, StoreId string, Ip string) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid":     UUID,
		"user_id":  UserId,
		"store_id": StoreId,
		"ip":       Ip,
		"exp":      time.Now().Add(time.Minute * 60),
	})
	token, err := tokenClaims.SignedString(secret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GeneratorIDCheckJWT(UserId string, condition Condition) (string, error) {
	t := time.Now()
	UUID, err := uuid.NewUUID()
	if err != nil {
		log.Debug("Get UUID Error", err)
		return "", err
	}
	config := viper.GetStringMapString("MOI")
	log.Debug("config", config)
	jsonCondition, err := json.Marshal(condition)
	if err != nil {
		log.Debug("Generator Json Error", err)
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":          "綠色便民專案",
		"orgId":        config["tax_id"],
		"apId":         config["ap_id"],
		"userId":       UserId[:12],
		"jobId":        config["job_id"],
		"opType":       "RW",
		"iss":          config["iss_key"],
		"iat":          t.Local().Add(-5 * time.Minute).Unix(), // start
		"aud":          t.Format("2006/01/02 15:04:05.000"),    // current time
		"exp":          t.Local().Add(5 * time.Minute).Unix(),  // expire
		"jti":          UUID,
		"conditionMap": string(jsonCondition),
	})
	jsonString, _ := json.Marshal(token)
	log.Debug("token", string(jsonString))
	prk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config["private_key"]))
	if err != nil {
		log.Debug("jwt Parse Private PRN Error ", err)
		return "", err
	}
	tokenString, err := token.SignedString(prk)
	if err != nil {
		log.Debug("Generator Validate ID jwt", err)
		return "", err
	}
	return tokenString, err
}

func ParseToken(token string) (*jwt.MapClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{},
		func(token *jwt.Token) (i interface{}, e error) {
			jsonString, _ := json.Marshal(secret)
			log.Debug("secret", string(jsonString))
			return secret, nil
		})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.MapClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}
	return nil, err
}

func validateJWT(token string) (status bool) {
	//algorithm := jwt.HmacSha256(secret)
	//validate := algorithm.Validate(token)
	//if validate != nil {
	//	panic(validate)
	//}
	//
	//loadedClaims, err := algorithm.Decode(token)
	//if err != nil {
	//	panic(err)
	//}
	//
	//role, err := loadedClaims.Get("Role")
	//if err != nil {
	//	panic(err)
	//}
	//
	//roleString, ok := role.(string)
	//if !ok {
	//	panic(err)
	//}

	//status = strings.Compare(roleString, "Admin") == 0

	return false
}
