package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jiharal/s1gu/example/myapp/api"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func HandlerUserLogin(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()

	var param api.UserLoginParam

	err := ParseBodyData(ctx, r, &param)
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "vehicle/create/param"),
			http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	_, token, aErr := userService.Login(ctx, param)
	if aErr != nil {
		return nil, aErr
	}

	data := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: token,
	}

	return data, nil
}
func HandlerUserRegister(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()

	var param api.UserDataParam

	err := ParseBodyData(ctx, r, &param)
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "vehicle/create/param"),
			http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	return userService.Register(ctx, param)
}

func HandlerUserList(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	return userService.List(ctx)
}

func HandlerUserDetail(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.FromString(vars["id"])
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "vehicle/detail"),
			http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	return userService.Detail(ctx, id)
}

func HandlerUserCreate(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()

	var param api.UserDataParam

	err := ParseBodyData(ctx, r, &param)
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "vehicle/create/param"),
			http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	param.By = api.GetContextRequesterID(ctx)
	return userService.Create(ctx, param)
}

func HandlerUserUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.FromString(vars["id"])
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "vehicle/update"),
			http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	var param api.UserDataParam
	err = ParseBodyData(ctx, r, &param)
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "vehicle/update/param"),
			http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	param.ID = id
	param.By = api.GetContextRequesterID(ctx)
	return nil, userService.Update(ctx, param)
}

func HandlerUserDelete(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.FromString(vars["id"])
	if err != nil {
		return nil, api.NewError(errors.Wrap(err, "vehicle/delete"),
			http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	return nil, userService.Delete(ctx, id)
}
