package dbPgs

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainU "github.com/Sskrill/TaskGyberNaty/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	userDB := NewUserDB(db)

	mock.ExpectExec("INSERT INTO users").
		WithArgs("test_user", "test_password").
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = userDB.CreateUser(context.Background(), domainU.User{
		UserName: "test_user",
		Password: "test_password",
	})
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	userDB := NewUserDB(db)

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(1, "test_user")
	mock.ExpectQuery("SELECT id,username FROM users WHERE username=\\$1 AND password=\\$2").
		WithArgs("test_user", "test_password").
		WillReturnRows(rows)

	user, err := userDB.GetUser(context.Background(), "test_password", "test_user")
	assert.NoError(t, err)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "test_user", user.UserName)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAllUsers(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userDB := NewUserDB(db)

	rows := sqlmock.NewRows([]string{"username"}).
		AddRow("test_user1").
		AddRow("test_user2")

	mock.ExpectQuery("SELECT username FROM users").
		WillReturnRows(rows)

	users, err := userDB.GetAllUsers(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "test_user1", users[0].UserName)
	assert.Equal(t, "test_user2", users[1].UserName)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
func TestGetUserNameById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	userDB := NewUserDB(db)

	mock.ExpectQuery("SELECT username FROM users WHERE id=\\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow("test_user"))
	userName, err := userDB.GetUserNameById(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "test_user", userName)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
