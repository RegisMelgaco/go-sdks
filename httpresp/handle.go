package httpresp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Handle(handler func(http.ResponseWriter, *http.Request) Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := handler(w, r)

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
		w.WriteHeader(resp.status)
	}
}
