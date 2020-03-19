package main

import (
	"demo-backend/handler"
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

func main() {
	router := mux.NewRouter()
	log.Println("-----------------Starting router----------------")
	router.HandleFunc("/form", handler.ConfigHandler).Methods("POST")
	router.HandleFunc("/data", handler.PersonHandler)

	log.Fatal(http.ListenAndServe(":3000", router))
}
