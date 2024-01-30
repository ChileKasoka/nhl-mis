package main

import (
	"log"
	"os"

	api "github.com/ChileKasoka/nhl-mis/api"
	"github.com/ChileKasoka/nhl-mis/db"
	sqlc "github.com/ChileKasoka/nhl-mis/db/sqlc"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("env")

	portString := os.Getenv("PORT")

	dbCon, err := db.ConnectDb()
	if err != nil {
		log.Fatal("failed to connect to db")
	}

	store := sqlc.NewStore(dbCon)
	server := api.NewServer(store)

	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	err = server.Start(":" + portString)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
