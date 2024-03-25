package server

import (
	"github.com/gorilla/mux"
	"goods/api/handler/good"
	"goods/api/handler/goods"
	"log"
	"net/http"
)

func Run() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/good/create", good.Create).Methods("POST")
	router.HandleFunc("/good/update", good.Update).Methods("PATCH")
	router.HandleFunc("/good/remove", good.Remove).Methods("DELETE")
	router.HandleFunc("/good/reprioritiize", good.Reprioritiize).Methods("PATCH")

	router.HandleFunc("/goods/list", goods.List).Methods("GET")

	log.Fatal(http.ListenAndServe(":8081", router))
}
