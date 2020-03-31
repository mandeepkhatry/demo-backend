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
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var eng engine.Engine
var port string

func init() {
	eng.DBName = "db"
	eng.Namespace = "npse"
	store := kvstore.NewBadgerFactory([]string{}, "./data/badger")
	eng.Store = store
	eng.ConnectDB()

	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("Error reading config file %s", err)
	// }
	// config_path := os.Getenv("config_path")
	// port = strings.Split(os.Getenv("http_target"), ":")[1]

	config_path := "./gateway/config.json"
	port = "3000"

	extra := make(map[string]string)
	extra["http_target"] = "localhost:3000"

	viper.SetConfigFile(config_path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	epoints := viper.Get("endpoints").([]interface{})
	if len(epoints) == 0 {
		epoints = append(epoints, utility.BuildFormPostConfig(extra)...)
		epoints = append(epoints, utility.BuildFormGetConfig(extra)...)
		epoints = append(epoints, utility.BuildDataGetConfig(extra)...)
		epoints = append(epoints, utility.BuildSearchConfig(extra)...)
		epoints = append(epoints, utility.BuildUpdateConfig(extra)...)
		epoints = append(epoints, utility.BuildDeleteConfig(extra)...)
		epoints = append(epoints, utility.BuildDataFetchConfig(extra)...)
		epoints = append(epoints, utility.BuildFormsGetConfig(extra)...)
		viper.Set("endpoints", epoints)
		viper.WriteConfig()
	}

}

func main() {
	router := mux.NewRouter()
	log.Println("-----------------Starting router----------------")
	router.HandleFunc("/form", ConfigHandler).Methods("POST")
	router.HandleFunc("/form", FormsGetHandler).Methods("GET")
	router.HandleFunc("/form/{table}", ConfigGetHandler).Methods("GET")
	router.HandleFunc("/api/table/{table}", DataPostHandler).Methods("POST")
	router.HandleFunc("/api/table/{table}", DataGetHandler).Methods("GET")
	router.HandleFunc("/api/table/{table}/{_id}", DataFetchHandler).Methods(("GET"))
	router.HandleFunc("/api/table/{table}/{_id}", UpdateHandler).Methods("PATCH")
	router.HandleFunc("/api/table/{table}/{_id}", DeleteHandler).Methods("DELETE")
	router.HandleFunc("/api/query", SearchHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func BuildKrakendConfig(config map[string]interface{}) error {

	http_target := "localhost:3000"
	endpoints := viper.Get("endpoints").([]interface{})

	extra_config := make(map[string]string)
	extra_config["http_target"] = http_target

	//add series of endpoints required
	endpoints = append(endpoints, utility.BuildDataPostConfig(config, extra_config)...)
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
	var response response.FormResponse
	response.Status = "config added successfully"
	statusCode := http.StatusOK

	for _, v := range viper.Get("endpoints").([]interface{}) {
		if v.(map[string]interface{})["endpoint"].(string) == "/api/table/"+config["table"].(string) {
			response.Status = "config added previously"
			encoding.JsonEncode(w, response, statusCode)
			return
		}
	}

	err := BuildKrakendConfig(config)

	if err != nil {
		response.Status = "config added unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}

	data := make(map[string][]byte)
	form := make(map[string][]byte)

	for k, v := range config {
		valueInBytes, err := json.Marshal(v)
		if err != nil {
			response.Status = "data added unsuccessfully"
			statusCode = http.StatusBadRequest
			encoding.JsonEncode(w, response, statusCode)
			return
		}
		if k == "table" {
			k = "_table"
			form["_table"] = valueInBytes
		} else if k == "schema" {
			k = "_schema"
		}
		data[k] = valueInBytes
	}

	docIndex := "bank_" + config["table"].(string) + "_schema"

	err = eng.InsertSingleIndexDocument(config["table"].(string), data, docIndex, "document")
	if err != nil {
		response.Status = "config added unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}

	err = eng.InsertSingleIndexDocument("_form", form, "", "")
	if err != nil {
		response.Status = "config added unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}
	encoding.JsonEncode(w, response, statusCode)
	return
}

func ConfigGetHandler(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	var response response.QueryResponse
	response.Status = "schema fetched successfully"
	response.Results = []map[string]interface{}{}
	statusCode := http.StatusOK

	docIndex := "bank_" + parameters["table"] + "_schema"

	resultArray, err := eng.SearchSingleDocument(parameters["table"], docIndex, "document")
	if err != nil {
		response.Status = "schema fetched unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}
	result := make([]map[string]interface{}, 0)

	for _, v := range resultArray {
		var resultInBytes = make(map[string][]byte)

		err := json.Unmarshal(v, &resultInBytes)
		if err != nil {
			response.Status = "schema fetched unsuccessfully"
			statusCode = http.StatusBadRequest
			encoding.JsonEncode(w, response, statusCode)
			return
		}
		eachResult := make(map[string]interface{})
		for key, val := range resultInBytes {
			if key == "_table" {
				key = "table"
			} else if key == "_schema" {
				key = "schema"
			}
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

	err := eng.InsertDocument(params["table"], data, indices)

	if err != nil {
		response.Status = "data added unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}

	encoding.JsonEncode(w, response, statusCode)
	return

}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	var request map[string]interface{}
	json.NewDecoder(r.Body).Decode(&request)

	query := request["query"].(string)

	var response response.QueryResponse
	response.Status = "query successful"
	response.Results = []map[string]interface{}{}
	statusCode := http.StatusOK

	collection, postfixQuery, err := parser.ParseQuery(query)
	if err != nil {
		response.Status = "query unsuccessful"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}

	resultArray, err := eng.SearchDocument(collection, postfixQuery)
	if err != nil {
		response.Status = "query unsuccessful"
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
			response.Status = "query unsuccessful"
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

func DataGetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var response response.QueryResponse
	response.Status = "query successful"
	response.Results = []map[string]interface{}{}
	statusCode := http.StatusOK

	parameters := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(strings.Trim(v[0], " ")) != 0 {
			parameters[k] = v[0]
		}
	}
	if len(parameters) != 0 {
		fmt.Println("PARAMETERS : ", parameters)
		parameters["table"] = params["table"]

		query := utility.ConvertParamsToQuery(parameters)
		fmt.Println("QUERY : ", query)

		collection, postfixQuery, err := parser.ParseQuery(query)
		if err != nil {
			response.Status = "query unsuccessful"
			statusCode = http.StatusBadRequest
			encoding.JsonEncode(w, response, statusCode)
			return
		}

		resultArray, err := eng.SearchDocument(collection, postfixQuery)
		if err != nil {
			response.Status = "query unsuccessful"
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
				response.Status = "query unsuccessful"
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

	} else {
		resultArray, err := eng.FetchAllDocuments(params["table"])
		if err != nil {
			response.Status = "query unsuccessful"
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
				response.Status = "query unsuccessful"
				statusCode = http.StatusBadRequest
				encoding.JsonEncode(w, response, statusCode)
				return
			}
			eachResult := make(map[string]interface{})
			if _, schemaok := resultInBytes["_schema"]; !schemaok {
				for key, val := range resultInBytes {
					var eachValue interface{}
					json.Unmarshal(val, &eachValue)
					eachResult[key] = eachValue
				}
				result = append(result, eachResult)
			}
		}

		response.Results = result

	}

	encoding.JsonEncode(w, response, statusCode)
	return

}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {

	var request map[string]interface{}
	json.NewDecoder(r.Body).Decode(&request)
	params := mux.Vars(r)
	data := make(map[string][]byte)
	indices := make([]string, 0)
	var response response.DataResponse
	response.Status = "data updated successfully"
	statusCode := http.StatusOK

	for k, v := range request {
		valueInBytes, err := json.Marshal(v)
		if err != nil {
			response.Status = "data updated unsuccessfully"
			statusCode = http.StatusBadRequest
			encoding.JsonEncode(w, response, statusCode)
			return
		}
		data[k] = valueInBytes
		indices = append(indices, k)
	}

	id, err := strconv.Atoi(params["_id"])
	if err != nil {
		response.Status = "data updated unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}
	err = eng.UpdateDocument(params["table"], data, id)

	if err != nil {
		response.Status = "data updated unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}

	encoding.JsonEncode(w, response, statusCode)
	return

}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var response response.DataResponse
	response.Status = "data deleted successfully"
	statusCode := http.StatusOK

	id, err := strconv.Atoi(params["_id"])
	if err != nil {
		response.Status = "data deleted unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}
	err = eng.DeleteDocument(params["table"], id)

	if err != nil {
		response.Status = "data deleted unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}

	encoding.JsonEncode(w, response, statusCode)
	return

}

func DataFetchHandler(w http.ResponseWriter, r *http.Request) {
	var response response.QueryResponse
	response.Status = "query successful"
	response.Results = []map[string]interface{}{}
	statusCode := http.StatusOK
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["_id"])

	if err != nil {
		response.Status = "query unsuccessful"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}
	resultArray, err := eng.SearchDocumentById(params["table"], id)
	if err != nil {
		response.Status = "query unsuccessful"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}
	var resultInBytes = make(map[string][]byte)

	err = json.Unmarshal(resultArray, &resultInBytes)
	delete(resultInBytes, "_indices")
	if err != nil {
		response.Status = "query unsuccessful"
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
	response.Results = append(response.Results, eachResult)
	encoding.JsonEncode(w, response, statusCode)
	return

}

func FormsGetHandler(w http.ResponseWriter, r *http.Request) {

	var response response.FormsResponse
	result := make([]interface{}, 0)
	response.Status = "forms fetched successfully"
	response.Results = result
	statusCode := http.StatusOK

	resultArray, err := eng.SearchDocumentByPrefix("_form")
	if err != nil {
		fmt.Println("ERROR : ", err)
		response.Status = "schema fetched unsuccessfully"
		statusCode = http.StatusBadRequest
		encoding.JsonEncode(w, response, statusCode)
		return
	}

	for _, v := range resultArray {
		var resultInBytes = make(map[string][]byte)

		err := json.Unmarshal(v, &resultInBytes)
		if err != nil {
			response.Status = "schema fetched unsuccessfully"
			statusCode = http.StatusBadRequest
			encoding.JsonEncode(w, response, statusCode)
			return
		}
		for _, val := range resultInBytes {
			var eachValue interface{}
			json.Unmarshal(val, &eachValue)
			result = append(result, eachValue)
		}

	}

	response.Results = result
	encoding.JsonEncode(w, response, statusCode)
	return

}
