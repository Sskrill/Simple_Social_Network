package dbPgs

import (
	"context"
	"database/sql"
	domainA "github.com/Sskrill/TaskGyberNaty/internal/domain/article"
)

type Article struct {
	Db *sql.DB
}

func NewArticleDB(db *sql.DB) *Article { return &Article{Db: db} }
func (a *Article) CreateArticle(ctx context.Context, article domainA.Article, userName string) error {
	_, err := a.Db.Exec("INSERT INTO article (user_name,title,content) VALUES ($1,$2,$3)",
		userName, article.Title, article.Content)
	return err
}
func (a *Article) GetAllArticlesByName(ctx context.Context, userName string) ([]domainA.Article, error) {
	rows, err := a.Db.Query("SELECT title,content FROM article WHERE user_name=$1", userName)
	if err != nil {
		return nil, err
	}
	articles := make([]domainA.Article, 0)
	for rows.Next() {
		article := domainA.Article{}
		if err := rows.Scan(article.Title, article.Content); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, rows.Err()
}
