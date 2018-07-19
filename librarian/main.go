package main

import (
	"github.com/chauhanr/searchy/librarian/common"
	"net/http"
	"github.com/chauhanr/searchy/librarian/api"
)

func main(){
	common.Log("Adding api handlers.. ")
	http.HandleFunc("/api/index", api.IndexHandler)

	common.Log("Starting indexer")
	api.StartIndexSystem()

	common.Log("Starting Searchy Librarian server at port 9090")
	http.ListenAndServe(":9090", nil)
}
