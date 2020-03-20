package response

type QueryResponse struct {
	Status  string                   `json:"status"`
	Results []map[string]interface{} `json:"results"`
}
