package manager

import (
	"../message"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// would be changed to tarantool
// it's just simple example to show how all works
type Student struct {
	Name string `json:"name"`
}

type PostForm struct {
	Id      int     `json:"id"`
	Student Student `json:"student"`
}

type Group map[int]Student

var group Group = make(map[int]Student)


// server listening API
type ClientManager struct {
	// Raft IO
	RaftIn  <-chan message.RaftMessage
	RaftOut chan<- message.RaftMessage
}

// handlers for different types of messages
func Post(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)

	log.Println(string(bytes))

	var form PostForm
	if err := json.Unmarshal(bytes, &form); err != nil {
		log.Print("invalid body")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(
			`{ "error" : "invalid body" }`,
		))
		return
	}

	log.Println(form)

	if _, exists := group[form.Id]; exists {
		log.Printf("Student already exists: %s", form.Student)

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(
			`{ "error" : "already exists" }`,
		))
		return
	}

	group[form.Id] = form.Student

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(
		`{ "result" : "Student posted" }`,
	))
}

func Get(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path[1:]

	id, err := strconv.Atoi(path)
	if err != nil {
		log.Print("invalid path")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(
			`{ "error" : "invalid path" }`,
		))
		return
	}

	if student, exists := group[id]; !exists {
		log.Printf("Student not found: %d", id)

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(
			`{ "error" : "Student not found" }`,
		))
		return
	} else {
		w.WriteHeader(http.StatusOK)

		type GetResponse struct {
			Result Student `json:"result"`
		}
		resp, _ := json.Marshal(GetResponse{Result: student})
		w.Write(resp)
	}
}

func Put(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path[1:]
	id, err := strconv.Atoi(path)
	if err != nil {
		log.Print("invalid path")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(
			`{ "error" : "invalid path" }`,
		))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var gotStudent Student
	if err := decoder.Decode(&gotStudent); err != nil {
		log.Print("invalid body")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(
			`{ "error" : "invalid body" }`,
		))
		return
	}

	if _, exists := group[id]; !exists {
		log.Printf("gotStudent not found: %d", id)

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(
			`{ "error" : "gotStudent not found" }`,
		))
		return
	} else {
		group[id] = gotStudent
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(
			`{ "result" : "Student updated" }`,
		))
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path[1:]

	id, err := strconv.Atoi(path)
	if err != nil {
		log.Print("invalid path")

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(
			`{ "error" : "invalid path" }`,
		))
		return
	}

	if _, exists := group[id]; !exists {
		log.Printf("Student not found: %d", id)

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(
			`{ "error" : "Student not found" }`,
		))
		return
	} else {
		delete(group, id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(
			`{ "result" : "Student deleted" }`,
		))
	}
}


func (cm* ClientManager) ClientManagerProcessEntrie() {
	r := mux.NewRouter()

	r.Methods("POST")  .Path("/")    .HandlerFunc(Post)
	r.Methods("GET")   .Path("/{Id}").HandlerFunc(Get)
	r.Methods("PUT")   .Path("/{Id}").HandlerFunc(Put)
	r.Methods("DELETE").Path("/{Id}").HandlerFunc(Delete)

	if err := http.ListenAndServe(":80", r); err != nil {
		log.Fatal(err)
	}
}
