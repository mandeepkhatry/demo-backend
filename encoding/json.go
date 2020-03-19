package encoding

import (
	"encoding/json"
	"net/http"
)

func JsonEncode(w http.ResponseWriter, response interface{}, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	return
}
