package handler

import (
	"demo-backend/encoding"
	"demo-backend/response"
	"demo-backend/utility"
	"encoding/json"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func BuildKrakendConfig(config map[string]interface{}) error {

	http_target := os.Getenv("http_target")
	endpoints := viper.Get("endpoints").([]interface{})

	extra_config := make(map[string]string)
	extra_config["http_target"] = http_target

	//add series of endpoints required
	endpoints = append(endpoints, utility.BuildFormGetConfig(config, extra_config)...)
	endpoints = append(endpoints, utility.BuildTablePostConfig(config, extra_config)...)
	endpoints = append(endpoints, utility.BuildTableGetConfig(config, extra_config)...)
	endpoints = append(endpoints, utility.BuildUpdateConfig(config, extra_config)...)

	viper.Set("endpoints", endpoints)

	err := viper.WriteConfig()
	return err

}

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	var config map[string]interface{}
	json.NewDecoder(r.Body).Decode(&config)
	err := BuildKrakendConfig(config)

	var response response.FormResponse
	response.Status = "config added successfully"
	statusCode := http.StatusOK
	if err != nil {
		response.Status = "config added unsuccessful"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}
	encoding.JsonEncode(w, response, statusCode)
	return

}
