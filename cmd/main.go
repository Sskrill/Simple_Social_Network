package main

import (
	"fmt"
	"github.com/Sskrill/TaskGyberNaty/internal/api"
	dbPgs "github.com/Sskrill/TaskGyberNaty/internal/repository/postgres"
	srvc "github.com/Sskrill/TaskGyberNaty/internal/service/user"
	"github.com/Sskrill/TaskGyberNaty/package/connDb"
	"github.com/Sskrill/TaskGyberNaty/package/hasher"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := connDb.NewDbPg()
	if err != nil {
		log.Fatal(err)
	}
	//salt:=os.Getenv()
	passwordHasher := hasher.NewHasher(os.Getenv("Salt"))
	articlesDb := dbPgs.NewArticleDB(db)
	userDb := dbPgs.NewUserDB(db)
	tokenDb := dbPgs.NewTokenDB(db)
	userSevice := srvc.NewServiceUser(tokenDb, userDb, passwordHasher, articlesDb, []byte(os.Getenv("Secret")))
	handler := api.NewHandler(userSevice)
	server := &http.Server{Addr: fmt.Sprintf(":8080"), Handler: handler.CreateRouter()}

	fmt.Println("Server started ")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
