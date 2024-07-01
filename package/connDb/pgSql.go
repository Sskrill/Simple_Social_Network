package connDb

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func NewDbPg() (*sql.DB, error) {
	cfg := ConfigParamDb{}
	if err := godotenv.Load(); err != nil {
		log.Fatal("error env")
	}
	if err := envconfig.Process("db", &cfg); err != nil {
		log.Fatal("cant take env")
	}
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBname, cfg.Sslmode))

	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
