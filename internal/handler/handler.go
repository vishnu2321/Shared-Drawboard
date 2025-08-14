package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
	"github.com/shared-drawboard/internal/models"
	"github.com/shared-drawboard/internal/service"
	"github.com/shared-drawboard/internal/websocket"
	"github.com/shared-drawboard/pkg/auth"
	"github.com/shared-drawboard/pkg/helper"
	"github.com/shared-drawboard/pkg/logger"
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
	router.HandleFunc("/refresh", h.refreshTokenHandler).Methods("POST")

	router.PathPrefix("/drawboard/").Handler(
		http.StripPrefix("/drawboard/", http.FileServer(http.Dir("./web/drawboard"))),
	).Methods("GET")

	router.HandleFunc("/drawboard", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/drawboard/", http.StatusMovedPermanently)
	}).Methods("GET")

	wsManager := websocket.NewManager()
	go wsManager.Run()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		h.websocketHandler(w, r, wsManager)
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

	authExpiryAt := time.Now().Add(15 * time.Minute).Unix()
	authtoken, err := auth.CreateJWTToken(user.Email, authExpiryAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	refreshTokenDTO, err := h.Service.CreateSession(r.Context(), user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	expiresAt, err := strconv.ParseInt(refreshTokenDTO.ExpiresAt, 10, 64)
	if err != nil {
		http.Error(w, "Invalid expires at value", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    refreshTokenDTO.TokenHash,
		Path:     "/",
		Expires:  time.Unix(expiresAt, 0),
		HttpOnly: true,                    // prevent JS access
		Secure:   true,                    // send only over HTTPS
		SameSite: http.SameSiteStrictMode, // CSRF protection
	})

	response := map[string]interface{}{
		"message":        "Sign in successful.",
		"auth-token":     authtoken,
		"auth-expiry-at": authExpiryAt,
		"user":           user,
	}
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenDTO

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	authToken := req.AuthToken
	if authToken == "" {
		http.Error(w, "Empty JWT tokken", http.StatusBadRequest)
		return
	}

	rtoken, err := r.Cookie("refresh-token")
	if err != nil {
		http.Error(w, "error reading token", http.StatusInternalServerError)
	}

	tdto, err := h.Service.UpdateSession(r.Context(), models.RefreshTokenDTO{AuthToken: authToken, RefreshToken: rtoken.Value})
	if err != nil {
		http.Error(w, "error updating session", http.StatusInternalServerError)
	}

	refreshExpiresAt, err := strconv.ParseInt(tdto.RefreshExpriesAt, 10, 64)
	if err != nil {
		http.Error(w, "Invalid expires at value", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    tdto.RefreshToken,
		Path:     "/refresh",
		Expires:  time.Unix(refreshExpiresAt, 0),
		HttpOnly: true,                    // prevent JS access
		Secure:   true,                    // send only over HTTPS
		SameSite: http.SameSiteStrictMode, // CSRF protection
	})

	response := map[string]interface{}{
		"message":        "Sign in successful.",
		"auth-token":     tdto.AuthToken,
		"auth-expiry-at": tdto.AuthExpiresAt,
	}
	json.NewEncoder(w).Encode(response)
}

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) websocketHandler(w http.ResponseWriter, r *http.Request, manager *websocket.Manager) {

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	claims, err := auth.VerifyJWTToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func(conn *ws.Conn, exp time.Time) {
		duration := time.Until(exp)
		if duration > 0 {
			time.Sleep(duration)
		}
		// Notify client
		_ = conn.WriteJSON(map[string]string{
			"type": "TOKEN_EXPIRED",
		})
		conn.Close()
	}(conn, exp.Time)

	client := &websocket.Client{
		ID:   helper.GenerateUniqueID(),
		Conn: conn,
		Send: make(chan []byte),
	}

	manager.Register <- client

	go handleRead(client, manager)
	go handleWrite(client, manager)
}

func handleRead(client *websocket.Client, manager *websocket.Manager) {
	defer func() {
		manager.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			logger.Error("Read error: %s", err)
			break
		}
		manager.Broadcast <- message
	}
}

func handleWrite(client *websocket.Client, _ *websocket.Manager) {
	for message := range client.Send {
		err := client.Conn.WriteMessage(ws.TextMessage, message)
		if err != nil {
			logger.Error("Write error: %s", err)
			break
		}
	}
	client.Conn.Close()
}
