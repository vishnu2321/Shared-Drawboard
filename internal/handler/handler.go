package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shared-drawboard/internal/middleware"
	"github.com/shared-drawboard/internal/models"
	"github.com/shared-drawboard/internal/service"
	"github.com/shared-drawboard/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Router  *mux.Router
	Service *service.Service
}

func New() (*Handler, error) {
	router := Router()
	service, err := service.New()
	if err != nil {
		return nil, err
	}

	h := &Handler{
		Router:  router,
		Service: service,
	}

	router.PathPrefix("/login/").Handler(
		http.StripPrefix("/login/", http.FileServer(http.Dir("./web/login"))),
	).Methods("GET")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login/", http.StatusMovedPermanently)
	})

	router.HandleFunc("/signup", h.signUpUserHandler).Methods("POST")
	router.HandleFunc("/signin", h.signinUserHandler).Methods("POST")

	router.PathPrefix("/drawboard/").Handler(
		middleware.AuthMiddleware(http.StripPrefix("/drawboard/", http.FileServer(http.Dir("./web/drawboard")))),
	).Methods("GET")

	router.HandleFunc("/drawboard", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/drawboard/", http.StatusMovedPermanently)
	})

	return h, nil
}

func (h *Handler) signUpUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	id, err := h.Service.SaveUser(r.Context(), req)
	if err != nil {
		if err == service.ErrUserExists {
			http.Error(w, "User already exisits", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"id": id}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) signinUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	user, err := h.Service.GetUser(r.Context(), req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Email != req.Email {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateNewToken(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := map[string]interface{}{"message": "sign in successful", "token": token}
	json.NewEncoder(w).Encode(response)

}
