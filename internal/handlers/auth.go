package handlers

import (
	"errors"
	"net/http"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

func NewAuthHandler(objectType string) (AuthHandler, error) {
	switch objectType {
	case "consumer":
		return &ConsumerAuthHandler{}, nil
	case "authorities":
		return &AuthoritiesAuthHandler{}, nil
	default:
		return nil, errors.New("invalid auth handler type")
	}
}
