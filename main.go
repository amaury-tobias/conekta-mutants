package main

import (
	"context"
	"log"
	"time"

	"github.com/amaury-tobias/conekta-mutants/internal/api"
	"github.com/amaury-tobias/conekta-mutants/internal/database"
)

func main() {
	err := database.Setup(database.NewMongoSession())
	checkError(err)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer database.DBClient.Disconnect(ctx)

	app := api.Init()
	checkError(app.Listen(":5000"))
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
