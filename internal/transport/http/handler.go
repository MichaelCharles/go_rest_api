package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcaubrey/go_rest_api/internal/services/comment"
	"github.com/mcaubrey/go_rest_api/internal/services/user"
	"github.com/sirupsen/logrus"
)

// Handler - stores pointer to our comments service
type Handler struct {
	Router         *mux.Router
	CommentService *comment.Service
	UserService    *user.Service
}

// Response - an object to store responses from our API
type Response struct {
	Message string
	Error   string
}

// NewHandler - returns a pointer to a Handler
func NewHandler(commentService *comment.Service, userService *user.Service) *Handler {
	return &Handler{
		CommentService: commentService,
		UserService:    userService,
	}
}

// SetupRoutes - sets up all the routes for our application
func (h *Handler) SetupRoutes() {
	logrus.Info("Setting up routes...")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)

	h.Router.HandleFunc("/api/comment", JWTAuth(h.GetAllComments)).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.GetComment)).Methods("GET")
	h.Router.HandleFunc("/api/comment", JWTAuth(h.PostComment)).Methods("POST")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.UpdateComment)).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", JWTAuth(h.DeleteComment)).Methods("DELETE")

	h.Router.HandleFunc("/api/user", h.GetAllUsers).Methods("GET")
	h.Router.HandleFunc("/api/login_user", h.LoginUser).Methods("POST")
	h.Router.HandleFunc("/api/register_user", h.RegisterUser).Methods("POST")

	h.Router.HandleFunc("/api/health", healthCheck)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := sendOKResponse(w, Response{Message: "Johnny 5 is alive!"}); err != nil {
		panic(err)
	}
}

func sendOKResponse(w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{
		Message: message,
		Error:   err.Error(),
	}); err != nil {
		panic(err)
	}
}
