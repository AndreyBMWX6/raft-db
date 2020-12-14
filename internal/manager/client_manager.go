package manager

import (
	"../message"
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

// server listening API
type ClientManager struct {
	// Raft IO
	ClientIn  <-chan message.ClientMessage
	ClientOut chan<- message.ClientMessage
	Port string
}

// handlers for different types of messages
func (cm *ClientManager) Post(w http.ResponseWriter, r *http.Request) {
	request := &message.DBRequest{
		Owner:   nil,
		Type:    message.PostRequestType,
		Request: r,
	}

	msg := message.NewRawClientMessage(
		&message.BaseClientMessage{
			Owner   : nil,
			Dest    : nil,
			ReqType : message.PostRequestType,
		},
		request,
	)

	cm.ClientOut<-msg

	response := <-cm.ClientIn
	switch resp := response.(type) {
	case *message.ResponseClientMessage:
		if resp.Redirect == true {
			http.Redirect(w, r, resp.LeaderURL + r.URL.Path, http.StatusPermanentRedirect)
		} else {
			w.WriteHeader(resp.DBResponse.Status)
			w.Write(resp.DBResponse.Result)
		}
	default:
		log.Print("`ResponseClientMessage` expected, got another type")
	}
}

func (cm *ClientManager) Get(w http.ResponseWriter, r *http.Request) {
	request := &message.DBRequest{
		Owner:   nil,
		Type:    message.GetRequestType,
		Request: r,
	}

	msg := message.NewRawClientMessage(
		&message.BaseClientMessage{
			Owner   : nil,
			Dest    : nil,
			ReqType : message.GetRequestType,
		},
		request,
	)

	cm.ClientOut<-msg

	response := <-cm.ClientIn
	switch resp := response.(type) {
	case *message.ResponseClientMessage:
		w.WriteHeader(resp.DBResponse.Status)
		w.Write(resp.DBResponse.Result)
	default:
		log.Print("`ResponseClientMessage` expected, got another type")
	}
}

func (cm *ClientManager) Put(w http.ResponseWriter, r *http.Request) {
		request := &message.DBRequest{
			Owner:   nil,
			Type:    message.PutRequestType,
			Request: r,
		}

		msg := message.NewRawClientMessage(
			&message.BaseClientMessage{
				Owner   : nil,
				Dest    : nil,
				ReqType : message.PutRequestType,
			},
			request,
		)

		cm.ClientOut<-msg

		response := <-cm.ClientIn
		switch resp := response.(type) {
		case *message.ResponseClientMessage:
			if resp.Redirect == true {
				http.Redirect(w, r, resp.LeaderURL + r.URL.Path, http.StatusPermanentRedirect)
			} else {
				w.WriteHeader(resp.DBResponse.Status)
				w.Write(resp.DBResponse.Result)
			}
		default:
			log.Print("`ResponseClientMessage` expected, got another type")
		}
}

func (cm *ClientManager) Delete(w http.ResponseWriter, r *http.Request) {
		request := &message.DBRequest{
			Owner:   nil,
			Type:    message.DeleteRequestType,
			Request: r,
		}

		msg := message.NewRawClientMessage(
			&message.BaseClientMessage{
				Owner   :   nil,
				Dest    :    nil,
				ReqType : message.DeleteRequestType,
			},
			request,
		)

		cm.ClientOut<-msg

		response := <-cm.ClientIn
		switch resp := response.(type) {
		case *message.ResponseClientMessage:
			if resp.Redirect == true {
				http.Redirect(w, r, resp.LeaderURL + r.URL.Path, http.StatusPermanentRedirect)
			} else {
				w.WriteHeader(resp.DBResponse.Status)
				w.Write(resp.DBResponse.Result)
			}
		default:
			log.Print("`ResponseClientMessage` expected, got another type")
		}
}

func (cm* ClientManager) ProcessEntries() {
	r := mux.NewRouter()

	r.Methods("POST")  .Path("/")    .HandlerFunc(cm.Post)
	r.Methods("GET")   .Path("/{Id}").HandlerFunc(cm.Get)
	r.Methods("PUT")   .Path("/{Id}").HandlerFunc(cm.Put)
	r.Methods("DELETE").Path("/{Id}").HandlerFunc(cm.Delete)

	if err := http.ListenAndServe(":" + cm.Port, r); err != nil {
		log.Fatal(err)
	}
}
