package router

import (
	"net/http"
	"log"

	"github.com/gorilla/mux"
	"../config"
)

type Router struct {
	URLs   []string
	CurrId int
}

func NewRouter() *Router {
	cfg := config.NewRouterConfig()
	return &Router{
		URLs   :   cfg.URLs,
		CurrId :   0,
	}
}

func NewRouterRunAll() *Router {
	cfg := config.NewRouterRunAllConfig()
	return &Router{
		URLs   :   cfg.URLs,
		CurrId :   0,
	}
}

func (router *Router) GetURL() *string {
	currId := router.CurrId
	router.CurrId = (router.CurrId + 1) % len(router.URLs)
	return &router.URLs[currId]
}

func (router *Router) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Redirect(w, r, *router.GetURL() + r.URL.Path, http.StatusPermanentRedirect)
	}
}

func (router *Router) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Redirect(w, r, *router.GetURL() + r.URL.Path, http.StatusPermanentRedirect)
	}
}

func (router *Router) Put(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		http.Redirect(w, r, *router.GetURL() + r.URL.Path, http.StatusPermanentRedirect)
	}
}

func (router *Router) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		http.Redirect(w, r, *router.GetURL() + r.URL.Path, http.StatusPermanentRedirect)
	}
}

func (router *Router) RunRouter() {
	r := mux.NewRouter()

	r.Methods("POST")  .Path("/")    .HandlerFunc(router.Post)
	r.Methods("GET")   .Path("/{Id}").HandlerFunc(router.Get)
	r.Methods("PUT")   .Path("/{Id}").HandlerFunc(router.Put)
	r.Methods("DELETE").Path("/{Id}").HandlerFunc(router.Delete)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

