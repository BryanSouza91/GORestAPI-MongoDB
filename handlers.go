package main

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	dao = DAO{}
)

// AllUsersEndpoint will GET list of users
func AllUsersEndpoint(w http.ResponseWriter, r *http.Request) {
	users, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

// FindUserEndpoint will GET a users by its ID
func FindUserEndpoint(w http.ResponseWriter, r *http.Request, id string) {
	user, err := dao.FindByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

// CreateUserEndpoint will POST a new user
func CreateUserEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		respondWithError(w, http.StatusBadRequest, "Please send a request body")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	user.ID = primitive.NewObjectID()
	if err := dao.Insert(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

// UpdateUserEndpoint will PUT update an existing user
func UpdateUserEndpoint(w http.ResponseWriter, r *http.Request, id string) {
	user.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := dao.Update(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

// DeleteUserEndpoint will DELETE an existing user
func DeleteUserEndpoint(w http.ResponseWriter, r *http.Request, id string) {
	deletedUser, err := dao.Delete(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, deletedUser)
}
