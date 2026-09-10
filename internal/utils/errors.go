package utils

import (
	"log"
	"net/http"
)

func IntervalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internalServeError: %s path %s, %s", r.Method, r.URL.Path, err.Error())
	WriteJsonError(w, http.StatusInternalServerError, err.Error())
}
func BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad Request: %s  path %s", err, r.URL.Path)
	WriteJsonError(w, http.StatusBadRequest, err.Error())
}

func UnauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Unauthorized: %s  path %s", err, r.URL.Path)
	WriteJsonError(w, http.StatusUnauthorized, err.Error())
}

func ForbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Forbidden: %s  path %s", err, r.URL.Path)
	WriteJsonError(w, http.StatusForbidden, err.Error())
}

func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not Found: %s  path %s", err, r.URL.Path)
	WriteJsonError(w, http.StatusNotFound, err.Error())
}

func TooManyRequestsError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Too Many Requests: %s  path %s", err, r.URL.Path)
	WriteJsonError(w, http.StatusTooManyRequests, err.Error())
}
