package api

import (
	"context"
	domainA "github.com/Sskrill/TaskGyberNaty/internal/domain/article"
	domainU "github.com/Sskrill/TaskGyberNaty/internal/domain/user"
)

type UserService interface {
	SignIn(ctx context.Context, param domainU.AuthParam) (string, string, error)
	SignUp(ctx context.Context, param domainU.AuthParam) error
	ParseToken(ctx context.Context, token string) (int64, error)
	GetAllArticles(ctx context.Context) ([]domainU.UserArticles /* -Возможно нужно ставить указатель */, error)
	FindArticlesByToken(ctx context.Context, rToken string, article domainA.Article) error
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

// Нужно реализовать
