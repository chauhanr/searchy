package api

import (
	"net/http"
	"encoding/json"
	"log"
)

var pProcessCh chan tPayload

var tcGetCh chan tCallback

func StartIndexSystem(){
	pProcessCh = make(chan tPayload, 100)
	tcGetCh = make(chan tCallback, 20)
	go tIndexer(pProcessCh, tcGetCh)
}

func tIndexer(ch chan tPayload, callback chan tCallback){
	store := tCatalog{}
	for {
		select {
			 case msg := <- callback:
				 dc := store[msg.Token]
				 msg.Ch <-tcMsg{
				 	DC: dc,
				 	Token: msg.Token,
				 }
			 case pd := <-ch:
			 	dc, exists := store[pd.Token]
			 	if !exists{
			 		dc = documentCatalog{}
			 		store[pd.Token] = dc
				}
			 	doc, exists := dc[pd.DocID]
			 	if !exists{
			 		doc = &document{
			 			DocID: pd.DocID,
			 			Title: pd.Title,
			 			Indices: map[int]tIndices{},
					}
			 		dc[pd.DocID] = doc
				}
			 	tin := tIndex{
			 		Index: pd.Index,
			 		LIndex: pd.LIndex,
				}
			 	doc.Indices[tin.LIndex] = append(doc.Indices[tin.LIndex], tin)
			 	doc.Count++
		}
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"code": 405, "msg": "Method Not allowed"}`))
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var tp tPayload
	decoder.Decode(&tp)
	log.Printf("Token received %#v\n", tp)
	pProcessCh <- tp

	w.Write([]byte(`{"code": 200, "msg": "Tokens are being added to index."}`))
}