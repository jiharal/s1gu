package router

import (
	"encoding/json"
	"net/http"

	"github.com/satori/go.uuid"

	"github.com/AuthScureDevelopment/authscure-go/example/MyApp/model"
	"github.com/gorilla/mux"
)

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := model.GetAllUser(r.Context(), DbPool)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string][]model.UserModel{"user": users})
}

func GetOneUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := model.GetOneUser(r.Context(), DbPool, id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]model.UserModel{"user": user})
}

func InsertUser(w http.ResponseWriter, r *http.Request) {
	var user model.UserModel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	id, err := user.Insert(r.Context(), DbPool)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user.ID = id

	respondWithJSON(w, http.StatusCreated, map[string]model.UserModel{"user": user})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var user model.UserModel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	user.ID = id

	if err := user.Update(r.Context(), DbPool); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"user": "success"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := model.DeleteUser(r.Context(), DbPool, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"user": "success"})
}
