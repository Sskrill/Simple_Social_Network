package api

import (
	"encoding/json"
	"fmt"
	domainErr "github.com/Sskrill/TaskGyberNaty/internal/domain/errors"
	domainU "github.com/Sskrill/TaskGyberNaty/internal/domain/user"
	"io/ioutil"
	"net/http"
)

func (h *Handler) refreshTokens(w http.ResponseWriter, r *http.Request) {
	cookieToken, err := r.Cookie("refresh-token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	aToken, rToken, err := h.userS.RefreshTokens(r.Context(), cookieToken.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	resp, err := json.Marshal(map[string]string{"token": aToken})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	cookie := &http.Cookie{
		Name:     "refresh-token",
		Value:    fmt.Sprintf("'%s'", rToken),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	dataReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	aParam := domainU.AuthParam{}
	if err := json.Unmarshal(dataReq, &aParam); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	if err = aParam.IsValid(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	aToken, rToken, err := h.userS.SignIn(r.Context(), aParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	resp, err := json.Marshal(map[string]string{"token": aToken})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})
		w.Write(resp)
		return
	}
	cookie := &http.Cookie{
		Name:     "refresh-token",
		Value:    fmt.Sprintf("'%s'", rToken),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)

}
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	dataReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	aParam := domainU.AuthParam{}
	if err := json.Unmarshal(dataReq, &aParam); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	if err = aParam.IsValid(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	err = h.userS.SignUp(r.Context(), aParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	w.WriteHeader(http.StatusOK)
}
