package httpresp

import (
	"net/http"
)

type Response struct {
	status  int
	payload any
	err     error
}

func OK(payload any) Response {
	return Response{
		status:  http.StatusOK,
		payload: payload,
	}
}

func Created(payload any) Response {
	return Response{
		status:  http.StatusCreated,
		payload: payload,
	}
}
