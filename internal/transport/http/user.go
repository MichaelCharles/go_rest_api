package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mcaubrey/go_rest_api/internal/services/user"
)

// GetUser - this will get a user by ID
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UINT from ID", err)
		return
	}

	user, err := h.UserService.GetUser(uint(i))
	if err != nil {
		sendErrorResponse(w, "Error getting user by ID", err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

// PostUser - creates a new user
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user user.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, "Enable to parse request body", err)
		return
	}

	user, err := h.UserService.RegisterUser(user)
	if err != nil {
		sendErrorResponse(w, "Error when posting new user.", err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

// PostUser - creates a new user
func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user user.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, "Enable to parse request body", err)
		return
	}

	user, err := h.UserService.LoginUser(user)
	if err != nil {
		sendErrorResponse(w, "Error when posting new user.", err)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

// GetAllUsers - gets all comments from comment service
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		sendErrorResponse(w, "Unable to get all users", err)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		panic(err)
	}
}
