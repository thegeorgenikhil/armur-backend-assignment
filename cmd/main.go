package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/thegeorgenikhil/armur-backend-assignment/internal/database"
	"github.com/thegeorgenikhil/armur-backend-assignment/internal/routes"
)

const port = ":8080"

func main() {
	database.InitDB("./data.db")
	fmt.Printf("Server started at port %s\n", port)

	srv := &http.Server{
		Addr:    port,
		Handler: routes.Routes(),
	}

	log.Fatal(srv.ListenAndServe())
}
