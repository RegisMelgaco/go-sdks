package handler

import (
	"encoding/json"
	"net/http"

	"github.com/regismelgaco/go-sdks/auth/auth/entity"
	v1 "github.com/regismelgaco/go-sdks/auth/auth/gateway/http/handler/v1"
	"github.com/regismelgaco/go-sdks/auth/auth/usecase"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/httpresp"
)

type Handler struct {
	u usecase.Usecase
}

func New(u usecase.Usecase) Handler {
	return Handler{u}
}

func (h Handler) PostUser(r *http.Request) httpresp.Response {
	var input v1.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		err = erring.Wrap(err).Build()

		return httpresp.BadRequest(err)
	}

	u, err := h.u.CreateUser(r.Context(), input.ToEntity())
	if err != nil {
		return httpresp.Error(err)
	}

	return httpresp.Created(v1.ToUserOutput(u))
}

func (h Handler) Login(r *http.Request) httpresp.Response {
	var input v1.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		err = erring.Wrap(err).Build()

		return httpresp.BadRequest(err)
	}

	t, err := h.u.Login(r.Context(), input.ToLoginEntity())
	if err != nil {
		return httpresp.Error(err)
	}

	return httpresp.OK(v1.ToLoginOutput(t))
}

func (h Handler) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			httpresp.BadRequest(entity.ErrMissingAuthorization).Handle(w)

			return
		}

		err := h.u.IsAuthorized(r.Context(), entity.Token(auth))
		if err != nil {
			httpresp.Error(err).Handle(w)

			return
		}

		next.ServeHTTP(w, r)
	})
}
