package api

import (
	"context"
	domainA "github.com/Sskrill/TaskGyberNaty/internal/domain/article"
	domainU "github.com/Sskrill/TaskGyberNaty/internal/domain/user"
	"github.com/gorilla/mux"
	"net/http"
)

type UserService interface {
	SignIn(ctx context.Context, param domainU.AuthParam) (string, string, error)
	SignUp(ctx context.Context, param domainU.AuthParam) error
	ParseToken(ctx context.Context, token string) (int64, error)
	GetAllArticles(ctx context.Context) (*[]domainU.UserArticles /* -Возможно нужно ставить указатель */, error)
	CraeteArticlesByToken(ctx context.Context, rToken string, article domainA.Article) error
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type Handler struct {
	userS UserService
}

func NewHandler(service UserService) *Handler { return &Handler{userS: service} }

func (h *Handler) CreateRouter() *mux.Router {
	router := mux.NewRouter()
	auth := router.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/refresh", h.refreshTokens).Methods(http.MethodGet)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodGet)
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
	}
	articles := router.PathPrefix("/article").Subrouter()
	{
		articles.Use(h.authMiddleware)
		articles.HandleFunc("", h.createArticle).Methods(http.MethodPost)
		articles.HandleFunc("/all", h.getAllArticles).Methods(http.MethodGet)
	}
	return router
}
