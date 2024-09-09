package ref

import (
	"awesomeProject/internal/http-server/handlers/tools"
	"awesomeProject/internal/http-server/services"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	Access  string `json:"access" validate:"required"`
	Refresh string `json:"refresh" validate:"required"`
}
type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewRefreshHandler(service services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if tools.ErrorHandler(err, http.StatusInternalServerError, w, r) {
			return
		}
		if err := validator.New().Struct(req); tools.ErrorHandler(err, http.StatusBadRequest, w, r) {
			return
		}
		access, refresh, err := service.Refresh.RefreshToken(req.Access, req.Refresh, tools.ReadUserIP(r))
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
