package dbPgs

import (
	"context"
	"database/sql"
	domainU "github.com/Sskrill/TaskGyberNaty/internal/domain/user"
)

type User struct {
	Db *sql.DB
}

func NewUserDB(db *sql.DB) *User { return &User{Db: db} }
func (u *User) CreateUser(ctx context.Context, user domainU.User) error {

	_, err := u.Db.Exec("INSERT INTO users (username,password) VALUES ($1,$2)", user.UserName, user.Password)
	return err
}
func (u *User) GetUser(ctx context.Context, password, userName string) (domainU.User, error) {
	user := domainU.User{}
	err := u.Db.QueryRow("SELECT id,username FROM users WHERE username=$1 AND password=$2", userName, password).
		Scan(&user.Id, &user.UserName)
	return user, err
}
func (u *User) GetAllUsers(ctx context.Context) ([]domainU.UserArticles, error) {
	rows, err := u.Db.Query("SELECT user_name FROM users")
	if err != nil {
		return nil, err
	}
	userArticles := make([]domainU.UserArticles, 0)
	for rows.Next() {
		userArtice := domainU.UserArticles{}
		if err = rows.Scan(&userArtice.UserName); err != nil {
			return nil, err
		}
		userArticles = append(userArticles, userArtice)
	}
	return userArticles, rows.Err()
}
func (u *User) GetUserNameById(ctx context.Context, id int) (string, error) {
	var name string
	err := u.Db.QueryRow("SELECT username FROM users WHERE id=$1", id).Scan(&name)
	return name, err
}
