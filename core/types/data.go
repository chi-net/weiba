package types

// for data.json

type ImportedGICAuthData struct {
	LastUpdate int            `json:"lastupdate"`
	Data       []gicGroupData `json:"data"`
}

type gicGroupData struct {
	ID   int64    `json:"i"`
	Data []string `json:"d"`
}
