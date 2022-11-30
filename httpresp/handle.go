package httpresp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/logger"
	"go.uber.org/zap"
)

func Handle(handler func(*http.Request) Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(r).Handle(w, r)
	}
}

func (res Response) Handle(w http.ResponseWriter, req *http.Request) {
	// write body
	err := json.NewEncoder(w).Encode(res.payload)
	if err != nil {
		//TODO log as json with Uberzap logger
		err = erring.Wrap(err).Describe("failed to encode response body")

		fmt.Println(err)
	}

	// write header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.status)

	// log error
	if res.err != nil {
		lvl := zap.ErrorLevel
		if res.status < http.StatusInternalServerError && res.status >= http.StatusBadRequest {
			lvl = zap.WarnLevel
		}

		logger := logger.FromContext(req.Context())

		var err erring.Err
		if errors.As(res.err, &err) {
			err.Log(logger, lvl)
		} else {
			logger.Log(lvl, res.err.Error())
		}
	}
}
