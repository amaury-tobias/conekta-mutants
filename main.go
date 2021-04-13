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
	mutantsService, err := SetupService(database.NewMongoSession())
	checkError(err)

	app := api.Init(mutantsService)
	checkError(app.Listen(":5000"))
}

func SetupService(session database.Session) (services.MutantsService, error) {
	db, err := database.NewMongoDatabase(session)
	checkError(err)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer db.Client().Disconnect(ctx)
	mutantsRepository := repository.NewMutantsRepository(db)
	return services.NewMutantsService(mutantsRepository), nil
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
