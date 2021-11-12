package dtos

type SearchFilter struct {
	Name    string   `json:"name"`
	Filters []Filter `json:"filters"`
}

type Filter struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}
