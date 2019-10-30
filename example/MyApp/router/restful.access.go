package router

import (
	"encoding/json"
	"net/http"

	"github.com/satori/go.uuid"

	"github.com/AuthScureDevelopment/authscure-go/example/MyApp/model"
	"github.com/gorilla/mux"
)

func GetAllAccess(w http.ResponseWriter, r *http.Request) {
	accessList, err := model.GetAllAccess(r.Context(), DbPool)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string][]model.AccessModel{"access": accessList})
}

func GetOneAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	access, err := model.GetOneAccess(r.Context(), DbPool, id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]model.AccessModel{"access": access})
}

func InsertAccess(w http.ResponseWriter, r *http.Request) {
	var access model.AccessModel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&access); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	id, err := access.Insert(r.Context(), DbPool)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	access.ID = id

	respondWithJSON(w, http.StatusCreated, map[string]model.AccessModel{"access": access})
}

func UpdateAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var access model.AccessModel
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&access); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	access.ID = id

	if err := access.Update(r.Context(), DbPool); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"access": "success"})
}

func DeleteAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := model.DeleteAccess(r.Context(), DbPool, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"access": "success"})
}
