package utils

import (
	"log"
	"net/http"
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

type PageError struct {
	Code int
	Err  error
	Msg  string
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

func (pe PageError) Error() string {
	return pe.Err.Error()
}

func (pe PageError) Status() int {
	return pe.Code
}

type Handler func(http.ResponseWriter, *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		switch e := err.(type) {
		case StatusError:
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		case PageError:
			log.Printf("HTTP %d - %s", e.Status(), e)
			RenderErrorTemplate(w, e.Msg, e.Status())
		default:
			log.Printf("Error: %s", e)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
