package api

import (
	"fmt"
	"bytes"
)

type tPayload struct{
	Token string `json:"token"`
	Title string `json:"title"`
	DocID string `json:"doc_id"`
	LIndex int `json:"line_index"`
	Index int  `json:"token_index"`
}

type tIndex struct{
	Index int
	LIndex int
}

type tIndices []tIndex

func (t *tIndex) String() string{
	return fmt.Sprintf("i: %d, li: %d", t.Index, t.LIndex)
}

type document struct{
	Count int
	DocID string
	Title string
	Indices map[int]tIndices
}

func (d *document) String() string{
	str := fmt.Sprintf("%s (%s): %d\n", d.Title, d.DocID, d.Count)
	var buffer bytes.Buffer
	for lin, tis := range d.Indices{
		var lBuffer bytes.Buffer
		for _, ti := range tis{
			lBuffer.WriteString(fmt.Sprintf("%s ", ti.String()))
		}
		buffer.WriteString(fmt.Sprintf("@%d -> %s\n", lin, lBuffer.String()))
	}
	return str+buffer.String()
}

type documentCatalog map[string]*document

func (d *documentCatalog) String() string{
	return fmt.Sprintf("%#v", d)
}

type tCatalog map[string]documentCatalog

func (t *tCatalog) String() string{
	return fmt.Sprintf("%#v", t)
}

type tCallback struct {
	Token string
	Ch 	  chan tcMsg
}

type tcMsg struct{
	Token string
	DC 	  documentCatalog
}

