package main

import (
	"github.com/chauhanr/searchy/concierge/common"
	"net/http"
	"github.com/chauhanr/searchy/concierge/api"
)

func main(){
	common.Log("Adding API handlers")
	http.HandleFunc("/api/feeder", api.FeedHandler)
	common.Log("Starting Feeder")
	api.StartFeederSystem()

	common.Log("Starting Searchy Concierge Service on port :8080")
	http.ListenAndServe(":8080", nil)
}
