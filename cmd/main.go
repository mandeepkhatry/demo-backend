package main

import (
	"demo-backend/encoding"
	"demo-backend/response"
	"demo-backend/server/engine"
	"demo-backend/server/engine/parser"
	"demo-backend/server/kvstore"
	"demo-backend/utility"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var eng engine.Engine

func init() {
	eng.DBName = "db"
	eng.Namespace = "npse"
	store := kvstore.NewBadgerFactory([]string{}, "./data/badger")
	eng.Store = store
	eng.ConnectDB()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error reading config file %s", err)
	}
	config_figpath := os.Getenv("krakend_config_path")
	viper.SetConfigFile(config_figpath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

}

func main() {
	router := mux.NewRouter()
	log.Println("-----------------Starting router----------------")
	router.HandleFunc("/form", ConfigHandler).Methods("POST")
	router.HandleFunc("/api/table/{name}", DataPostHandler).Methods("POST")
	router.HandleFunc("/api/query", DataGetHandler).Methods("POST")
	router.HandleFunc("/api/table/search", SearchHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func BuildKrakendConfig(config map[string]interface{}) error {

	http_target := os.Getenv("http_target")
	endpoints := viper.Get("endpoints").([]interface{})

	extra_config := make(map[string]string)
	extra_config["http_target"] = http_target

	//add series of endpoints required
	endpoints = append(endpoints, utility.BuildFormGetConfig(config, extra_config)...)
	endpoints = append(endpoints, utility.BuildDataPostConfig(config, extra_config)...)
	endpoints = append(endpoints, utility.BuildDataGetConfig(config, extra_config)...)
	endpoints = append(endpoints, utility.BuildUpdateConfig(config, extra_config)...)

	// for _, v := range endpoints {
	// 	// fmt.Println(v.(map[string]interface{})["endpoint"])
	// }

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
		response.Status = "config added unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}
	encoding.JsonEncode(w, response, statusCode)
	return

}

func DataPostHandler(w http.ResponseWriter, r *http.Request) {
	var request map[string]interface{}
	json.NewDecoder(r.Body).Decode(&request)
	params := mux.Vars(r)
	data := make(map[string][]byte)
	indices := make([]string, 0)
	var response response.DataResponse
	response.Status = "data added successfully"
	statusCode := http.StatusOK

	for k, v := range request {
		fmt.Println(k, v)
		valueInBytes, err := json.Marshal(v)
		if err != nil {
			response.Status = "data added unsuccessfully"
			statusCode = http.StatusBadRequest
			encoding.JsonEncode(w, response, statusCode)
			return
		}
		data[k] = valueInBytes
		indices = append(indices, k)
	}

	err := eng.InsertDocument(params["name"], data, indices)

	if err != nil {
		response.Status = "data added unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}

	encoding.JsonEncode(w, response, statusCode)
	return

}

func DataGetHandler(w http.ResponseWriter, r *http.Request) {
	var request map[string]interface{}
	json.NewDecoder(r.Body).Decode(&request)

	query := request["query"].(string)

	var response response.QueryResponse
	response.Status = "query successfully"
	response.Results = []map[string]interface{}{}
	statusCode := http.StatusOK

	collection, postfixQuery, err := parser.ParseQuery(query)
	if err != nil {
		response.Status = "query unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}

	resultArray, err := eng.SearchDocument(collection, postfixQuery)
	if err != nil {
		response.Status = "query unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}

	result := make([]map[string]interface{}, 0)

	for _, v := range resultArray {
		var resultInBytes = make(map[string][]byte)

		err := json.Unmarshal(v, &resultInBytes)
		delete(resultInBytes, "_indices")
		if err != nil {
			response.Status = "query unsuccessfully"
			statusCode = http.StatusBadRequest
			encoding.JsonEncode(w, response, statusCode)
			return
		}
		eachResult := make(map[string]interface{})
		for key, val := range resultInBytes {
			var eachValue interface{}
			json.Unmarshal(val, &eachValue)
			eachResult[key] = eachValue
		}
		result = append(result, eachResult)
	}

	response.Results = result
	encoding.JsonEncode(w, response, statusCode)
	return

}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	return
}
