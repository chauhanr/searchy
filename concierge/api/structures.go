package api

type payload struct{
	URL string `json:"url"`
	Title string `json:"title"`
}

type document struct {
	Doc string `json:"-"`
	Title string `json:"title"`
	DocID string `json:"docId"`
}

type token struct{
	Line string `json:"-"`
	Token string `json:"token"`
	Title string `json:"title"`
	DocID string `json:"doc_id"`
	LIndex int `json:"line_index"`
	Index int `json:"token_index"`
}

type dMsg struct{
	DocID string
	Ch    chan document
}

type lMsg struct {
	LIndex 		int
	DocID 		string
	Ch chan 	string
}

type lMeta struct {
	LIndex 		int
	DocID 		string
	Line 	 	string
}

type dAllMsg struct{
	Ch chan []document
}

