package httpresp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Handle(handler func(*http.Request) Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(r).Handle(w)
	}
}

func (resp Response) Handle(w http.ResponseWriter) {
	// log error
	if resp.err != nil {
		//TODO log as json with Uberzap logger

		fmt.Println(resp.err)
	}

	// write body
	err := json.NewEncoder(w).Encode(resp.payload)
	if err != nil {
		//TODO log as json with Uberzap logger

		fmt.Println(resp.err)
	}

	// write header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.status)
}
