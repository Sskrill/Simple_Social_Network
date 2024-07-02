package dbPgs

import (
	"context"
	"database/sql"
	domain "github.com/Sskrill/TaskGyberNaty/internal/domain/token"
	"strings"
)

type Token struct {
	Db *sql.DB
}

func NewTokenDB(db *sql.DB) *Token { return &Token{Db: db} }
func (t *Token) CreateToken(ctx context.Context, token domain.RefreshToken) error {
	_, err := t.Db.Exec("INSERT INTO refreshtokens (token,expires_at,user_id) VALUES ($1,$2,$3)",
		token.Token, token.ExpiresAt, token.UserID)
	return err
}
func (t *Token) GetToken(ctx context.Context, token string) (domain.RefreshToken, error) {
	refreshToken := domain.RefreshToken{}
	token = strings.ReplaceAll(token, "'", "") //-Убираем '' в начале и в конце
	err := t.Db.QueryRow("SELECT id,token,user_id,expires_at FROM refreshtokens WHERE token=$1", token).
		Scan(&refreshToken.Id, &refreshToken.Token, &refreshToken.UserID, &refreshToken.ExpiresAt)

	if err != nil {
		return refreshToken, err
	}

	err = t.DeleteTokenByUserId(ctx, refreshToken.UserID)

	return refreshToken, err
}
func (t *Token) GetUserIdByToken(ctx context.Context, token string) (int, error) {
	var id int
	token = strings.ReplaceAll(token, "'", "")
	err := t.Db.QueryRow("SELECT user_id FROM refreshtokens WHERE token=$1", token).Scan(&id)
	return id, err
}
func (t *Token) DeleteTokenByUserId(ctx context.Context, userId int) error {
	_, err := t.Db.Exec("DELETE FROM refreshtokens WHERE user_id=$1", userId)
	return err
}
