package dbPgs

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	domain "github.com/Sskrill/TaskGyberNaty/internal/domain/token"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	tokenDB := NewTokenDB(db)
	mock.ExpectExec("INSERT INTO refreshtokens").
		WithArgs("test_token", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = tokenDB.CreateToken(context.Background(), domain.RefreshToken{
		Token:     "test_token",
		ExpiresAt: time.Now(),
		UserID:    1,
	})
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	tokenDB := NewTokenDB(db)
	rows := sqlmock.NewRows([]string{"id", "token", "user_id", "expires_at"}).
		AddRow(1, "test_token", 1, time.Now())
	mock.ExpectQuery("SELECT id,token,user_id,expires_at FROM refreshtokens WHERE token=\\$1").
		WithArgs("test_token").
		WillReturnRows(rows)

	mock.ExpectExec("DELETE FROM refreshtokens WHERE user_id=\\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	token, err := tokenDB.GetToken(context.Background(), "test_token")
	assert.NoError(t, err)
	assert.Equal(t, 1, token.Id)
	assert.Equal(t, "test_token", token.Token)
	assert.Equal(t, 1, token.UserID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserIdByToken(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	tokenDB := NewTokenDB(db)
	mock.ExpectQuery("SELECT user_id FROM refreshtokens WHERE token=\\$1").
		WithArgs("test_token").
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
	userId, err := tokenDB.GetUserIdByToken(context.Background(), "test_token")
	assert.NoError(t, err)
	assert.Equal(t, 1, userId)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteTokenByUserId(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	tokenDB := NewTokenDB(db)

	mock.ExpectExec("DELETE FROM refreshtokens WHERE user_id=\\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = tokenDB.DeleteTokenByUserId(context.Background(), 1)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
