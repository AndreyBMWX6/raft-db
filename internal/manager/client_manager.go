package manager

import (
	"../message"
	//"encoding/binary"
	"fmt"
	//"log"
	"net/http"
)

// server listening API
type ClientManager struct {
	// Raft IO
	RaftIn  <-chan message.RaftMessage
	RaftOut chan<- message.RaftMessage
}

func Handler(w http.ResponseWriter, r *http.Request) {
	buff := make([]byte, r.ContentLength)
	r.Body.Read(buff)
	if buff == nil {
		fmt.Fprintf(w, "no data\n")
	}
	fmt.Fprintf(w, "Hello World!")
}

func (cm* ClientManager) ClientManagerProcessEntrie() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":80", nil)
}