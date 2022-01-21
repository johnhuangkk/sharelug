package TokenVo

type TokenParams struct {
	Platform        string `json:"Platform"`
	PlatformVersion string `json:"PlatformVersion"`
	PlatformDevice  string `json:"PlatformDevice"`
	FcmToken        string `json:"FcmToken"`
}

type CheckTokenParams struct {
	Uuid  string `json:"Uuid"`
	Token string `json:"Token"`
}

type ChangeUuidResponse struct {
	UUID    string `json:"UUID"`
}

type ChangeTokenResponse struct {
	UUID    string `json:"UUID"`
	Token   string `json:"Token"`
	Message int64  `json:"Message"`
}
