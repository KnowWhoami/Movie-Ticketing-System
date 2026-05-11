package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	jsonHelper "KnowWhoami/movie-ticketing/internal/json"
	userSvc "KnowWhoami/movie-ticketing/pkgs/service/user"
)

func (h *Handler) CreateUser(request *http.Request, writer http.ResponseWriter) {
	var inp userSvc.CreateUserInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}
	out, err := h.svc.user.CreateUser(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "CreateUser", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusCreated)
}

func (h *Handler) ListUsers(request *http.Request, writer http.ResponseWriter) {
	q := request.URL.Query()
	inp := userSvc.ListUsersInput{
		Name:     q.Get("name"),
		Email:    q.Get("email"),
		UserType: q.Get("user_type"),
	}
	out, err := h.svc.user.ListUsers(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "ListUsers", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) GetUserByID(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id == 0 {
		resp := ErrorResponse("invalid user id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}
	out, err := h.svc.user.GetUserByID(id)
	if err != nil {
		_ = h.logger.Log("handler", "GetUserByID", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusNotFound)
		jsonHelper.WriteResult(&resp, writer, http.StatusNotFound)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) UpdateUser(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id == 0 {
		resp := ErrorResponse("invalid user id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}
	var inp userSvc.UpdateUserInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}
	inp.ID = id
	out, err := h.svc.user.UpdateUser(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "UpdateUser", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) DeleteUser(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id == 0 {
		resp := ErrorResponse("invalid user id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}
	if err := h.svc.user.DeleteUser(id); err != nil {
		_ = h.logger.Log("handler", "DeleteUser", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(nil)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}
