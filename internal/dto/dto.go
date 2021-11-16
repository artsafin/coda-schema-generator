package dto

type Tables struct {
	Items []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"items"`
}

type TableColumns struct {
	TableID string
	Items   []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"items"`
}
