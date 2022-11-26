package httpresp

import (
	"errors"
	"net/http"

	"github.com/regismelgaco/go-sdks/erring"
)

type ErrResponse struct {
	Description string `json:"description"`
	//TODO error code
}

func newErrDescResponse(err error, status int, fallbackDesc string) Response {
	payload := ErrResponse{Description: fallbackDesc}

	var e erring.Err
	if errors.As(err, &e) {
		if e.Name != "" {
			payload.Description = e.Name
		} else if e.Description != "" {
			payload.Description = e.Description
		}
	}

	return Response{
		status:  status,
		payload: payload,
		err:     err,
	}
}

func BadRequest(err error) Response {
	return newErrDescResponse(err, http.StatusBadRequest, "bad request")
}

func NotFound(err error) Response {
	return newErrDescResponse(err, http.StatusNotFound, "not found")
}

func Unauthorized(err error) Response {
	return newErrDescResponse(err, http.StatusNotFound, "unauthorized")
}

func Internal(err error) Response {
	return Response{
		status:  http.StatusInternalServerError,
		payload: ErrResponse{Description: "internal server error"},
		err:     err,
	}
}

func Error(err error) Response {
	switch {
	case errors.Is(err, erring.ErrBadRequest):
		return BadRequest(err)
	case errors.Is(err, erring.ErrNotFound):
		return NotFound(err)
	case errors.Is(err, erring.ErrUnauthorized):
		return Unauthorized(err)
	}

	return Internal(err)
}
