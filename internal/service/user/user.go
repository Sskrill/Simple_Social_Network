package srvc

import (
	"context"
	"errors"
	"fmt"
	domainA "github.com/Sskrill/TaskGyberNaty/internal/domain/article"
	domainErr "github.com/Sskrill/TaskGyberNaty/internal/domain/errors"
	domainT "github.com/Sskrill/TaskGyberNaty/internal/domain/token"
	domainU "github.com/Sskrill/TaskGyberNaty/internal/domain/user"
	"github.com/golang-jwt/jwt"
	"math/rand"

	"regexp"
	"strconv"
	"time"
)

type ArticleRepo interface {
	CreateArticle(ctx context.Context, article domainA.Article, userName string) error
	GetAllArticlesByName(ctx context.Context, userName string) ([]domainA.Article, error)
}
type TokenRepo interface {
	GetToken(ctx context.Context, token string) (domainT.RefreshToken, error)
	CreateToken(ctx context.Context, token domainT.RefreshToken) error
	GetUserIdByToken(ctx context.Context, token string) (int, error)
	DeleteTokenByUserId(ctx context.Context, userId int) error
}
type UserRepo interface {
	GetUser(ctx context.Context, password, userName string) (domainU.User, error)
	CreateUser(ctx context.Context, user domainU.User) error
	GetUserNameById(ctx context.Context, id int) (string, error)
	GetAllUsers(ctx context.Context) ([]domainU.UserArticles, error)
}
type Hasher interface {
	Hash(str string) (string, error)
}
type ServiceUser struct {
	TRepo  TokenRepo
	URepo  UserRepo
	ARepo  ArticleRepo
	Hsh    Hasher
	Secret []byte
}

func isLatin(s string) bool {
	regex := regexp.MustCompile("^[a-zA-Z]+$")
	return regex.MatchString(s)
}
func validateArticle(input string) bool {

	regex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	return regex.MatchString(input)
}

func NewServiceUser(tRepo TokenRepo, uRepo UserRepo, hsh Hasher, aRepo ArticleRepo, secret []byte) *ServiceUser {
	return &ServiceUser{TRepo: tRepo, URepo: uRepo, Hsh: hsh, ARepo: aRepo, Secret: secret}
}
func (sU *ServiceUser) CraeteArticlesByToken(ctx context.Context, rToken string, article domainA.Article) error {
	if ok := validateArticle(article.Content); !ok {
		return errors.New("invalid content or title")
	}
	userId, err := sU.TRepo.GetUserIdByToken(ctx, rToken)
	if err != nil {

		return err
	}
	name, err := sU.URepo.GetUserNameById(ctx, userId)
	if err != nil {
		return err
	}
	err = sU.ARepo.CreateArticle(ctx, article, name)
	return err
}
func (sU *ServiceUser) GetAllArticles(ctx context.Context) ([]domainU.UserArticles /* -Возможно нужно ставить указатель */, error) {
	userArticels, err := sU.URepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, user := range userArticels {
		user.Articles, err = sU.ARepo.GetAllArticlesByName(ctx, user.UserName)
		if err != nil {
			return nil, err
		}
	}
	return userArticels, err // -Возможно нужно ставить указатель
}
func (sU *ServiceUser) SignIn(ctx context.Context, param domainU.AuthParam) (string, string, error) {
	password, err := sU.Hsh.Hash(param.Password)
	if err != nil {
		return "", "", err
	}
	user, err := sU.URepo.GetUser(ctx, password, param.UserName)
	if err != nil {
		return "", "", err
	}
	return sU.generateTokens(ctx, user.Id)
}
func (sU *ServiceUser) SignUp(ctx context.Context, param domainU.AuthParam) error {
	if !isLatin(param.UserName) {
		return domainErr.ErrorInvalidUsername
	}
	password, err := sU.Hsh.Hash(param.Password)
	if err != nil {

		return err
	}
	user := domainU.User{UserName: param.UserName, Password: password}
	return sU.URepo.CreateUser(ctx, user)
}
func (sU *ServiceUser) generateTokens(ctx context.Context, userId int) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userId)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	})

	accessToken, err := t.SignedString(sU.Secret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}
	err = sU.TRepo.DeleteTokenByUserId(ctx, userId)
	if err != nil {
		return "", "", err
	}
	if err := sU.TRepo.CreateToken(ctx, domainT.RefreshToken{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
func (sU *ServiceUser) ParseToken(ctx context.Context, token string) (int64, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return sU.Secret, nil
	})
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}
func (sU *ServiceUser) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := sU.TRepo.GetToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domainErr.ErrorRefreshTokenExpired
	}

	return sU.generateTokens(ctx, session.UserID)
}
