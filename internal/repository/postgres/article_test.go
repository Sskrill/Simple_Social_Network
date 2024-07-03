package dbPgs

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainA "github.com/Sskrill/TaskGyberNaty/internal/domain/article"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	articleDB := NewArticleDB(db)

	mock.ExpectExec("INSERT INTO article").
		WithArgs("test_user", "test_title", "test_content").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = articleDB.CreateArticle(context.Background(), domainA.Article{
		Title:   "test_title",
		Content: "test_content",
	}, "test_user")
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAllArticlesByName(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	articleDB := NewArticleDB(db)
	rows := sqlmock.NewRows([]string{"title", "content"}).
		AddRow("test_title1", "test_content1").
		AddRow("test_title2", "test_content2")

	mock.ExpectQuery("SELECT title,content FROM article WHERE user_name=\\$1").
		WithArgs("test_user").
		WillReturnRows(rows)
	articles, err := articleDB.GetAllArticlesByName(context.Background(), "test_user")
	assert.NoError(t, err)
	assert.Len(t, articles, 2)
	assert.Equal(t, "test_title1", articles[0].Title)
	assert.Equal(t, "test_content1", articles[0].Content)
	assert.Equal(t, "test_title2", articles[1].Title)
	assert.Equal(t, "test_content2", articles[1].Content)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
