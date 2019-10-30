
			package router

			import (
				"encoding/json"
				"net/http"

				"github.com/satori/go.uuid"

				"github.com/AuthScureDevelopment/authscure-go/example/MyApp/model"
				"github.com/gorilla/mux"
			)

			func GetAllJihar(w http.ResponseWriter, r *http.Request) {
				jiharList, err := model.GetAllJihar(r.Context(), DbPool)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				
				respondWithJSON(w, http.StatusOK, map[string][]model.JiharModel{"jihar": jiharList})
			}

			func GetOneJihar(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				id, err := uuid.FromString(vars["id"])
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				jihar, err := model.GetOneJihar(r.Context(), DbPool, id)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}

				respondWithJSON(w, http.StatusOK, map[string]model.JiharModel{"jihar": jihar})
			}

			func InsertJihar(w http.ResponseWriter, r *http.Request) {
				var jihar model.JiharModel
				decoder := json.NewDecoder(r.Body)
				if err := decoder.Decode(&jihar); err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				defer r.Body.Close()

				id, err := jihar.Insert(r.Context(), DbPool)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				jihar.ID = id

				respondWithJSON(w, http.StatusCreated, map[string]model.JiharModel{"jihar": jihar})
			}

			func UpdateJihar(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				id, err := uuid.FromString(vars["id"])
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}

				var jihar model.JiharModel
				decoder := json.NewDecoder(r.Body)
				if err := decoder.Decode(&jihar); err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				defer r.Body.Close()

				jihar.ID = id

				if err := jihar.Update(r.Context(), DbPool); err != nil {
					respondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}

				respondWithJSON(w, http.StatusOK, map[string]string{"jihar": "success"})
			}

			func DeleteJihar(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				id, err := uuid.FromString(vars["id"])
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}

				if err := model.DeleteJihar(r.Context(), DbPool, id); err != nil {
					respondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}

				respondWithJSON(w, http.StatusOK, map[string]string{"jihar": "success"})
			}
			