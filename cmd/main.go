package main

import (
	"demo-backend/encoding"
	"demo-backend/response"
	"demo-backend/utility"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error reading config file %s", err)
	}
	config_figpath := os.Getenv("krakend_config_path")
	viper.SetConfigFile(config_figpath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

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

func PersonHandler(w http.ResponseWriter, r *http.Request) {
	var config map[string]interface{}
	json.NewDecoder(r.Body).Decode(&config)

	var response response.FormResponse
	response.Status = "config added successfully"
	statusCode := http.StatusOK
	encoding.JsonEncode(w, response, statusCode)
	return

}
func main() {
	router := mux.NewRouter()
	log.Println("-----------------Starting router----------------")
	router.HandleFunc("/form", ConfigHandler).Methods("POST")
	router.HandleFunc("/data", PersonHandler)

	log.Fatal(http.ListenAndServe(":3000", router))
}
