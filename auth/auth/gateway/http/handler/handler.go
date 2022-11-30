package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/regismelgaco/go-sdks/auth/auth/entity"
	"github.com/regismelgaco/go-sdks/auth/auth/gateway/encryptor"
	v1 "github.com/regismelgaco/go-sdks/auth/auth/gateway/http/handler/v1"
	"github.com/regismelgaco/go-sdks/auth/auth/gateway/postgres/repository"
	"github.com/regismelgaco/go-sdks/auth/auth/usecase"
	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/httpresp"
)

type Handler struct {
	u usecase.Usecase
}

func NewHandler(p *pgxpool.Pool, jwtSecret string) Handler {
	encry := encryptor.NewEncryptor(jwtSecret)
	repo := repository.NewUserRepository(p)
	u := usecase.NewUsecase(encry, repo)

	return Handler{u}
}

func (h Handler) SetupRoutes(r chi.Router) {
	r.Post("/signup", httpresp.Handle(h.PostUser))
	r.Post("/login", httpresp.Handle(h.Login))
}

func (h Handler) PostUser(r *http.Request) httpresp.Response {
	var input v1.UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		err = erring.Wrap(err)

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
		err = erring.Wrap(err)

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
			httpresp.BadRequest(entity.ErrMissingAuthorization).Handle(w, r)

			return
		}

		claims, err := h.u.IsAuthorized(r.Context(), entity.Token(auth))
		if err != nil {
			httpresp.Error(err).Handle(w, r)

			return
		}

		ctx := context.WithValue(r.Context(), claimsKey{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type claimsKey struct{}

func ClaimsFromContext(ctx context.Context) (entity.TokenClaims, error) {
	c, ok := ctx.Value(claimsKey{}).(entity.TokenClaims)
	if !ok {
		return entity.TokenClaims{}, erring.Wrap(entity.ErrMissingClaimsCtx)
	}

	return c, nil
}
