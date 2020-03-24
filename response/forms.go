package response

type FormsResponse struct {
	Status  string        `json:"status"`
	Results []interface{} `json:"results"`
}
