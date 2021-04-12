package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/amaury-tobias/conekta-mutants/internal/api"
	"github.com/amaury-tobias/conekta-mutants/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	name     string
	method   string
	route    string
	want     string
	body     []byte
	wantCode int
}

func TestTestRoute(t *testing.T) {
	tests := []Test{
		{
			name:     "mutant valid",
			method:   "POST",
			route:    "/mutant",
			wantCode: fiber.StatusOK,
			want:     "OK",
			body:     []byte(`{"dna":["ATGCGA","CAGTGC","TTATGT","AGAAGG","CCCCTA","TCACTG"]}`),
		},
		{
			name:     "mutant invalid: no human",
			method:   "POST",
			route:    "/mutant",
			wantCode: fiber.StatusBadRequest,
			want:     "la secuencia de ADN contiene bases no validas",
			body: []byte(`{
				"dna":[
					"ATGCGA",
					"CAGTGC",
					"TTGMTT",
					"AGAAGG",
					"CACCTA",
					"CACCTA"
					]
				}`),
		},
		{
			name:     "mutant invalid: bad sequencies size",
			method:   "POST",
			route:    "/mutant",
			wantCode: fiber.StatusBadRequest,
			want:     "largo de secuencia invalido",
			body: []byte(`{
				"dna":[
					"ATGCGA",
					"CAGTGC",
					"TTGMTT",
					"AGAAGG"
					]
				}`),
		},
		{
			name:     "no mutant",
			method:   "POST",
			route:    "/mutant",
			wantCode: fiber.StatusForbidden,
			want:     "Forbidden",
			body: []byte(`{
				"dna":[
					"ATGCGA",
					"CAGTGC",
					"TTGTTT",
					"AGAAGG",
					"CACCTA",
					"TCACTG"
					]
				}`),
		},
		{
			name:     "Not found",
			method:   "GET",
			route:    "/i-dont-exists",
			want:     "Not Found",
			wantCode: fiber.StatusNotFound,
		},
	}

	err := database.Setup(database.NewMockSession())
	checkError(err)

	app := api.Init()

	for _, test := range tests {
		req, _ := http.NewRequest(
			test.method,
			test.route,
			bytes.NewBuffer(test.body),
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)

		assert.Nilf(t, err, test.name)
		assert.Equalf(t, test.wantCode, res.StatusCode, test.name)

		body, err := ioutil.ReadAll(res.Body)

		assert.Nilf(t, err, test.name)
		assert.Equalf(t, test.want, string(body), test.name)
	}
}
