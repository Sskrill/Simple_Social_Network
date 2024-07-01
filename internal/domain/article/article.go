package domain

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type Article struct {
	Title   string `json:"Title" validate:"required,gte=3,lte=100"`
	Content string `json:"Content" validate:"required"`
}

func (a Article) IsValid() error { return validate.Struct(a) }
