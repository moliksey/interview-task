package tools

import (
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Error struct {
	Error string `json:"message"`
}

func ErrorHandler(err error, status int, w http.ResponseWriter, r *http.Request) bool {
	if err != nil {
		var errText string
		switch status {
		case http.StatusBadRequest:
			errText = "invalid request"
		case http.StatusUnauthorized:
			errText = "invalid data"
		case http.StatusInternalServerError:
			errText = "error with decoding"
		default:
			errText = "smth went wrong"
		}
		logrus.Info(errText)
		render.JSON(w, r, Error{
			Error: errText,
		})
		w.WriteHeader(status)
		w.Header().Add("Content-Type", "application/json")
		return true
	}
	return false
}
func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
