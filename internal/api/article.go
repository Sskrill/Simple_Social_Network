package api

import (
	"encoding/json"
	"fmt"
	domainA "github.com/Sskrill/TaskGyberNaty/internal/domain/article"
	domainErr "github.com/Sskrill/TaskGyberNaty/internal/domain/errors"
	"io/ioutil"
	"net/http"
)

func (h *Handler) createArticle(w http.ResponseWriter, r *http.Request) {

	cookieRToken, err := r.Cookie("refresh-token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	resp, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	article := domainA.Article{}
	if err = json.Unmarshal(resp, &article); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}
	if err = article.IsValid(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}

	if err = h.userS.CraeteArticlesByToken(r.Context(), cookieRToken.Value, article); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})

		w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h *Handler) getAllArticles(w http.ResponseWriter, r *http.Request) {
	userArticles, err := h.userS.GetAllArticles(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})
		fmt.Println(2)
		w.Write(resp)
		return
	}
	resp, err := json.Marshal(userArticles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(domainErr.ErrorResponse{Message: err.Error()})
		fmt.Println(3)
		w.Write(resp)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
	return
}
