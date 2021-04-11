package main

import (
	"log"

	"github.com/amaury-tobias/conekta-mutants/internal/api"
	"github.com/amaury-tobias/conekta-mutants/internal/database"
)

func main() {
	err := database.Init()
	checkError(err)
	defer database.Disconnect()

	app := api.Init()
	checkError(app.Listen(":5000"))
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
