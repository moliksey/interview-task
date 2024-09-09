package auth

import (
	"awesomeProject/internal/http-server/handlers/tools"
	"awesomeProject/internal/http-server/services"
	"github.com/go-chi/render"
	"net/http"
)

type Request struct {
	Guid string `json:"guid" validate:"required"`
}
type Response struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func NewAUTHHandler(service services.Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if tools.ErrorHandler(err, http.StatusInternalServerError, w, r) {
			return
		}
		guids, present := r.Form["guid"]
		if !present || len(guids) != 1 {
			tools.ErrorHandler(err, http.StatusBadRequest, w, r)
			return
		}
		access, refresh, err := service.Authorization.AuthUser(guids[0], tools.ReadUserIP(r))
		if tools.ErrorHandler(err, http.StatusUnauthorized, w, r) {
			return
		}
		render.JSON(w, r, Response{
			AccessToken:  access,
			RefreshToken: refresh,
		})
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
	}
}
