package reg

import (
	"awesomeProject/internal/http-server/handlers/tools"
	"awesomeProject/internal/http-server/services"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Request struct {
	Email string `json:"email" validate:"required,email"`
}
type Response struct {
	Guid string `json:"guid"`
}

func NewRegHandler(service services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if tools.ErrorHandler(err, http.StatusInternalServerError, w, r) {
			return
		}
		if err := validator.New().Struct(req); tools.ErrorHandler(err, http.StatusBadRequest, w, r) {
			return
		}
		var guid string
		guid, err = service.AddUser(req.Email)
		if tools.ErrorHandler(err, http.StatusUnauthorized, w, r) {
			return
		}
		render.JSON(w, r, Response{
			Guid: guid,
		})
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
	}
}
