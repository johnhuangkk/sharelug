package IPOSTVO

/**
匯出DB POSTDATA
會用到
*/
type IPostJson struct {
	Country string `json:"country"`
	City    []City
}

type City struct {
	Country string `json:"country"`
	Name    string `json:"city"`
	Zip     string `json:"zip"`
	Address []Address
}

type Address struct {
	Name     string `json:"address"`
	Location string `json:"Location"`
	Id       string `json:"adm_id"`
	Alias    string `json:"adm_alias"`
}
