package handler

import (
	"demo-backend/encoding"
	"demo-backend/response"
	"encoding/json"
	"net/http"
)

func PersonHandler(w http.ResponseWriter, r *http.Request) {
	var config map[string]interface{}
	json.NewDecoder(r.Body).Decode(&config)

	var response response.FormResponse
	response.Status = "config added successfully"
	statusCode := http.StatusOK
	encoding.JsonEncode(w, response, statusCode)
	return

}
