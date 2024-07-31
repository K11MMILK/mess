package api

import (
	"github.com/gorilla/mux"
)

func NewRouter(handler *Handler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/messages", handler.PostMessage).Methods("POST")
	router.HandleFunc("/stats", handler.GetStats).Methods("GET")
	return router
}
