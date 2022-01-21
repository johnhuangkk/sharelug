package Response

type TWIDResponse struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
}

type TWIDVerify struct {
	Response TWIDVerifyResponse `json:"responseData"`
	HttpCode string             `json:"httpCode"`
	HttpMsg  string             `json:"httpMessage"`
	RespMsg  string             `json:"rdMessage"`
	RespCode string             `json:"rdCode"`
}

type TWIDVerifyResponse struct {
	CheckIDCardApply string `json:"checkIdCardApply"`
}
