package main

import (
	"context"
	"log"
	"time"

	"github.com/amaury-tobias/conekta-mutants/internal/api"
	"github.com/amaury-tobias/conekta-mutants/internal/database"
	"github.com/amaury-tobias/conekta-mutants/internal/repository"
	"github.com/amaury-tobias/conekta-mutants/internal/services"
)

func main() {
	db, err := database.NewMongoDatabase(database.NewMongoSession())
	checkError(err)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer db.Client().Disconnect(ctx)
	mutantsRepository := repository.NewMutantsRepository(db)
	mutantsService := services.NewMutantsService(mutantsRepository)

	app := api.Init(mutantsService)
	checkError(app.Listen(":5000"))
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
