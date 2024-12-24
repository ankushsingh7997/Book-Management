package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ankush/bookstore/pkg/models"
	"github.com/ankush/bookstore/pkg/utils"
	"github.com/gorilla/mux"
)

func CreateUser(res http.ResponseWriter, req *http.Request) {
	user := &models.User{}
	utils.ParseBody(req, user)

	newUser, err := user.CreateUser()
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(newUser)
}

func GetAllUsers(res http.ResponseWriter, req *http.Request) {
	users, err := models.GetAllUsers()
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(map[string]string{
			"error": "Failed to fetch users",
		})
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(users)
}

func GetUserByID(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userId"]

	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(map[string]string{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := models.GetUserByID(uint(id))
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(map[string]string{
			"error": "User not found",
		})
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(user)
}

func UpdateUser(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userId"]

	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(map[string]string{
			"error": "Invalid user ID",
		})
		return
	}

	existingUser, err := models.GetUserByID(uint(id))
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(map[string]string{
			"error": "User not found",
		})
		return
	}

	updateData := &models.User{}
	utils.ParseBody(req, updateData)

	// Update fields if provided
	if updateData.Name != "" {
		existingUser.Name = updateData.Name
	}
	if updateData.Email != "" {
		existingUser.Email = updateData.Email
	}
	if updateData.Phone != "" {
		existingUser.Phone = updateData.Phone
	}
	if updateData.Password != "" {
		existingUser.Password = updateData.Password
	}

	if err := existingUser.UpdateUser(); err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(map[string]string{
			"error": "Failed to update user",
		})
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(existingUser)
}

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["userId"]

	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(map[string]string{
			"error": "Invalid user ID",
		})
		return
	}

	if err := models.DeleteUser(uint(id)); err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(map[string]string{
			"error": "Failed to delete user",
		})
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]string{
		"message": "User deleted successfully",
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Body)

	var loginReq LoginRequest
	utils.ParseBody(req, &loginReq)
	fmt.Println(loginReq)

	user, err := models.GetUserByEmail(loginReq.Email)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(map[string]string{
			"error": "invalid credentials",
		})
		return
	}
	token, err := user.Login(loginReq.Password)
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	http.SetCookie(res, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Login successful",
		"user_id": user.ID,
		"email":   user.Email,
	})

}
