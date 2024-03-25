package httputils

import (
	"encoding/json"
	"errors"
	"goods/pkg/validate"
	"log"
	"net/http"
)

func ResponseBody(w http.ResponseWriter, response interface{}, headerStatus int) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(headerStatus)
	w.Write(jsonResponse)
}

func ParseBody(w http.ResponseWriter, r *http.Request, body interface{}) {
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
		log.Println("Error decoding JSON body:", err)
		return
	}

	resValid := validate.Struct(body)
	if len(resValid) > 0 {
		http.Error(w, "Failed to parse JSON body", http.StatusBadRequest)
		log.Println("Error decoding JSON body:", errors.New(resValid.Error()))
		return
	}
}
