package apiserver

import (
	"errors"
	"net/http"
)


type errorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

var (
	ErrParseURL            = errors.New("url parse error")
	ErrUnprocessableEntity = errors.New("data parse error")
	ErrNotFound            = errors.New("not found")
	ErrInternalServerError = errors.New("internal error")
	ErrEmptyData           = errors.New("empty data")
	ErrInvalidURL          = errors.New("invalid url")

	ErrResp500 = errorResponse{
		Code: http.StatusInternalServerError,
		Error: ErrInternalServerError.Error(),
	}

	
)
