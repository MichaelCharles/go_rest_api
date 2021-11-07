package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mcaubrey/go_rest_api/internal/services/comment"
)

// Handler - stores pointer to our comments service
type Handler struct {
	Router         *mux.Router
	CommentService *comment.Service
}

// Response - an object to store responses from our API
type Response struct {
	Message string
	Error   string
}

// NewHandler - returns a pointer to a Handler
func NewHandler(commentService *comment.Service) *Handler {
	return &Handler{
		CommentService: commentService,
	}
}

// SetupRoutes - sets up all the routes for our application
func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up routes...")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/comment", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment", h.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", h.DeleteComment).Methods("DELETE")

	h.Router.HandleFunc("/api/health", healthCheck)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := sendOKResponse(w, Response{Message: "Johnny 5 is alive!"}); err != nil {
		panic(err)
	}
}

// GetComment - this will get a comment by ID
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UINT from ID", err)
		return
	}

	comment, err := h.CommentService.GetComment(uint(i))
	if err != nil {
		sendErrorResponse(w, "Error getting comment by ID", err)
		return
	}

	if err := json.NewEncoder(w).Encode(comment); err != nil {
		panic(err)
	}
}

// GetAllComments - gets all comments from comment service
func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.CommentService.GetAllComments()
	if err != nil {
		sendErrorResponse(w, "Unable to get all comments", err)
		return
	}

	if err := json.NewEncoder(w).Encode(comments); err != nil {
		panic(err)
	}
}

// PostComment - creates a new comment
func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var comment comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "Enable to parse request body", err)
		return
	}

	comment, err := h.CommentService.PostComment(comment)
	if err != nil {
		sendErrorResponse(w, "Error when posting new comment.", err)
		return
	}

	if err := json.NewEncoder(w).Encode(comment); err != nil {
		panic(err)
	}
}

// UpdateComment - updates a comment by ID
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UINT from ID", err)
		return
	}

	var comment comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "Unable to parse request body", err)
		return
	}

	comment, err = h.CommentService.UpdateComment(uint(i), comment)
	if err != nil {
		sendErrorResponse(w, "Error when updating comment", err)
		return
	}

	if err := json.NewEncoder(w).Encode(comment); err != nil {
		panic(err)
	}
}

// DeleteComment - deletes a comment by ID
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UINT from ID", err)
		return
	}

	err = h.CommentService.DeleteComment(uint(i))
	if err != nil {
		sendErrorResponse(w, "Error when updating comment", err)
		return
	}

	if err = sendOKResponse(w, Response{Message: "Successfully deleted comment"}); err != nil {
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
