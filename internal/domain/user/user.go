package domain

import (
	domainA "github.com/Sskrill/TaskGyberNaty/internal/domain/article"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type User struct {
	Id       int
	UserName string
	Password string
}

type AuthParam struct {
	UserName string `json:"username" validate:"required,gte=3"`
	Password string `json:"password" validate:"required,gte=8"`
}

type UserArticles struct {
	Articles []domainA.Article `json:"articles"`
	UserName string            `json:"username"`
}

func (aP AuthParam) IsValid() error { return validate.Struct(aP) }
