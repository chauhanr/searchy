package api

import (
	"net/http"
	"fmt"
	"github.com/chauhanr/searchy/concierge/common"
	"crypto/sha1"
	"strings"
	"time"
	"io/ioutil"
	"encoding/json"
)

var done chan bool
// dGetCh - used to retrieve single document from store.
var dGetCh chan dMsg
// lGetCh - use to retrieve single line form the store
var lGetCh chan lMsg
// lStoreCh used to put line in store
var lStoreCh chan lMeta

// iAddCh is used to add token to index
var iAddCh chan token
// dStoreCh is used to put document in store
var dStoreCh chan document
// dProcessCh it is to break document to tokens and process them.
var dProcessCh chan document

// dGetAllCh use to get all teh documents from store.
var dGetAllCh chan dAllMsg
// process payload channel and start indexing process
var pProcessCh chan payload


/**
	The StartFeederSystem or service will start the feeder routines that can handle multiple documents and payload.
	This processing of the document is very parrallel in nature and we can see multiple documents being processed.
*/

func StartFeederSystem(){
	 done = make(chan bool)
	 dGetCh = make(chan dMsg, 8)
	 dGetAllCh = make(chan dAllMsg)

	 iAddCh = make(chan token, 8)
	 pProcessCh = make(chan payload, 8)

	 dStoreCh = make(chan document, 8)
	 dProcessCh = make(chan document, 8)
	 lGetCh = make(chan lMsg)
	 lStoreCh = make(chan lMeta, 8)

	 for i:=0; i<4; i++{
		go indexAdder(iAddCh, done)
		go docProcessor(pProcessCh, dStoreCh, dProcessCh, done)
		go indexProcessor(dProcessCh, lStoreCh, iAddCh, done)
		go docStore(dStoreCh, dGetCh, dGetAllCh, done)
		go lineStore(lStoreCh, lGetCh, done)
	 }
}

func indexAdder(ch chan token, done chan bool){
	for {
		select {
			case tok := <-ch:
				fmt.Println("adding to librarian: ", tok.Token)
			case <-done:
				common.Log("Exiting Index Addr.")
				return
		}
	}
}

func indexProcessor(ch chan document, lStoreCh chan lMeta, iAddCh chan token, done chan bool){
	for {
		select{
			case doc := <-ch:
				docLines := strings.Split(doc.Doc, "\n")
				lin := 0
				for _, line := range docLines {
					if strings.TrimSpace(line) == "" {
						continue
					}
					lStoreCh <- lMeta{
						LIndex: lin,
						Line:   line,
						DocID:  doc.DocID,
					}
					index := 0
					words := strings.Fields(line)
					for _, word := range words {
						if tok, valid := common.SimplifyToken(word); valid {
							iAddCh <- token{
								Token:  tok,
								LIndex: lin,
								Line:   line,
								Index:  index,
								DocID:  doc.DocID,
								Title:  doc.Title,
							}
							index++
						}
					}
					lin++
				}
			case <-done:
				common.Log("Exiting indexProcessor")
				return
		}

	}
}

func docStore(add chan document, get chan dMsg, dGetAllCh chan dAllMsg, done chan bool){
	store := map[string]document{}
	for{
		select{
			case doc := <-add:
				store[doc.DocID] = doc
			case m := <-get:
				m.Ch <- store[m.DocID]
			case ch := <-dGetAllCh:
				docs := []document{}
				for _, doc := range store{
					docs = append(docs, doc)
				}
				ch.Ch <- docs
			case <-done:
				common.Log("Exiting docStore")
				return
		}
	}
}

func docProcessor(in chan payload, dStoreCh chan document, dProcessCh chan document, done chan bool){
	for {
		select {
			case newDoc := <-in:
				var err error

				doc := ""
				if doc, err = getFile(newDoc.URL); err != nil{
					common.Warn(err.Error())
					continue
				}
				titleId := getTitleHash(newDoc.Title)
				msg := document {
					Doc: doc,
					DocID: titleId,
					Title: newDoc.Title,
				}

				dStoreCh <- msg
				dProcessCh <- msg
			case <-done:
				common.Log("Existing docProcessor.")
				return
		}
	}
}

func getTitleHash(title string) string{
	hash := sha1.New()
	title = strings.ToLower(title)

	str := fmt.Sprintf("%s-%s", time.Now(), title)
	hash.Write([]byte(str))
	hByte := hash.Sum(nil)
	return fmt.Sprintf("%x", hByte)
}

func getFile(URL string) (string, error){
	var res *http.Response
	var err error
	if res, err = http.Get(URL); err != nil{
		errMsg := fmt.Errorf("Unable to retireve server URL: %s, \n Error: %s", URL, err.Error())
		return "", errMsg
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil{
		errMsg := fmt.Errorf("Error while reading response: URL:%s.\n Status Code: %d, \nError: %s", URL, res.StatusCode, err.Error())
		return "", errMsg
	}
	return string(body), nil
}

func lineStore(ch chan lMeta, callBackChannel chan lMsg, done chan bool){
	store := map[string]string{}
	for {
		select{
			case line := <-ch:
				id := fmt.Sprintf("%s-%d", line.DocID, line.LIndex)
				store[id] = line.Line
			case ch := <-callBackChannel:
				line := ""
				id := fmt.Sprintf("%s-%d", ch.DocID, ch.LIndex)
				if l, exists := store[id]; exists{
					line = l
				}
				ch.Ch <- line
			case <-done:
				common.Log("Existing docStore")
				return
		}
	}
}

/**
	This is the front end api that is exposed and intercepts all calls to the /api/feeder urls.
*/
func FeedHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {
		ch := make(chan []document)
		dGetAllCh <- dAllMsg{ch}
		docs := <-ch
		close(ch)
		if serializedPayload, err := json.Marshal(docs); err == nil {
			w.Write(serializedPayload)
		}else{
			common.Warn("Unable to serialize the documents searched.")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"code": 500, "msg": "Error occurred while trying to retrieve documents.}`))
		}
		return
	}else if r.Method != http.MethodPost{
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`{"code": 405, "msg": "Method not allowed.}`))
			return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var newDoc payload
	decoder.Decode(&newDoc)
	pProcessCh <- newDoc

	w.Write([]byte(`{"code": 200, "msg": "Request is being processed."}`))
}
