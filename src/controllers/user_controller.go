package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go_chat/src/models"
	"go_chat/src/services"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	service *services.UserService
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Picture  string `json:"picture"`
	Password string `json:"password" validate:"required"`
}

type AuthRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{service}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerRequest RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&registerRequest)

	if err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.service.FindByEmail(context.Background(), registerRequest.Email)

	if err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if user != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"status":  "error",
			"message": "Email is already in use",
		})

		return
	}

	user = &models.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: registerRequest.Password,
	}

	newUser, err := c.service.Register(context.Background(), user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(map[string]any{
			"status":  "error",
			"message": "There was an error trying to register the user!",
		})
	}

	token, err := c.service.GenerateToken(newUser)
	if err != nil {
		http.Error(w, "There was an error trying to generate the token", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"user":  newUser,
		"token": token,
	})
}

func (c *UserController) AuthUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authRequest AuthRequest

	err := json.NewDecoder(r.Body).Decode(&authRequest)

	if err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.service.FindByEmail(context.Background(), authRequest.Email)

	if err != nil {
		log.Fatal()
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]any{
			"status":  "error",
			"message": fmt.Sprintf("User with email '%s' not found!", authRequest.Email),
		})

		return
	}

	if !c.service.MatchPassword(user, authRequest.Password) {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(map[string]any{
			"status":  "error",
			"message": "Email and password does not match!",
		})

		return
	}

	token, err := c.service.GenerateToken(user)

	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"user":  user,
		"token": token,
	})
}

func (c *UserController) SearchUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	search := strings.Trim(r.URL.Query().Get("search"), "")

	if search == "" {
		w.WriteHeader(http.StatusUnauthorized)

		json.NewEncoder(w).Encode(map[string]any{
			"status":  "error",
			"message": "Search cannot be empty!",
		})

		return
	}

	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{
				"email": primitive.Regex{
					Pattern: fmt.Sprintf("^%s+[a-zA-Z]+@[a-z]+\\.([a-z]+)?$", search),
					Options: "i",
				},
			},
		},
	}

	users, err := c.service.SearchUsers(context.Background(), filter)

	if err != nil {
		http.Error(w, "Error decoding JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if users == nil {
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]any{
			"message": "No user found!",
		})

		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"users": users,
	})
}
