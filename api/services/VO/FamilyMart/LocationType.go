package FamilyMart



type City struct {
	Country string `json:"country"`
	Districts    []District
}

type District struct {
	Name    string `json:"city"`
	Zip     string `json:"zip"`
	Address []Address
}

type Address struct {
	Name     string `json:"address"`
	Location string `json:"location"`
	StoreId       string `json:"id"`
	StoreName    string `json:"alias"`
}

