
			package router

			import (
				"encoding/json"
				"net/http"

				"github.com/satori/go.uuid"

				"github.com/AuthScureDevelopment/authscure-go/example/MyApp/model"
				"github.com/gorilla/mux"
			)

			func GetAllHello(w http.ResponseWriter, r *http.Request) {
				helloList, err := model.GetAllHello(r.Context(), DbPool)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				
				respondWithJSON(w, http.StatusOK, map[string][]model.HelloModel{"hello": helloList})
			}

			func GetOneHello(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				id, err := uuid.FromString(vars["id"])
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				hello, err := model.GetOneHello(r.Context(), DbPool, id)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}

				respondWithJSON(w, http.StatusOK, map[string]model.HelloModel{"hello": hello})
			}

			func InsertHello(w http.ResponseWriter, r *http.Request) {
				var hello model.HelloModel
				decoder := json.NewDecoder(r.Body)
				if err := decoder.Decode(&hello); err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				defer r.Body.Close()

				id, err := hello.Insert(r.Context(), DbPool)
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				hello.ID = id

				respondWithJSON(w, http.StatusCreated, map[string]model.HelloModel{"hello": hello})
			}

			func UpdateHello(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				id, err := uuid.FromString(vars["id"])
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}

				var hello model.HelloModel
				decoder := json.NewDecoder(r.Body)
				if err := decoder.Decode(&hello); err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}
				defer r.Body.Close()

				hello.ID = id

				if err := hello.Update(r.Context(), DbPool); err != nil {
					respondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}

				respondWithJSON(w, http.StatusOK, map[string]string{"hello": "success"})
			}

			func DeleteHello(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				id, err := uuid.FromString(vars["id"])
				if err != nil {
					respondWithError(w, http.StatusBadRequest, err.Error())
					return
				}

				if err := model.DeleteHello(r.Context(), DbPool, id); err != nil {
					respondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}

				respondWithJSON(w, http.StatusOK, map[string]string{"hello": "success"})
			}
			