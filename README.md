# conekta-mutants

## Desafíos:

- Nivel 1

  - Analizar secuencia de ADN y encontrar si esta es mutante
  - Metodo con firma `boolean IsMutant(String[] dna)`

- Nivel 2 - API REST

  - endpoint `/mutant` recibe mediante `POST`

  ```json
  { "dna”:["ATGCGA","CAGTGC","TTATGT","AGAAGG","CCCCTA","TCACTG"] }
  ```

  - Si es mutante devuelve `200-OK` si es no mutante `403-Forbidden`

- Nivel 3 - Base de Datos y Ratio
  - endpoint `/stats` mediante `GET` responde
  ```json
  { "count_mutant_dna":40, "count_human_dna”:100: "ratio":0.4 }
  ```
  - Guardado en base de datos

## Tecnologías Usadas

- GO ([Fiber](https://docs.gofiber.io/))
- MongoDB
- Docker
- Traefic

## Setup

### Requisitos

- Go 1.16
- MongoDB > 4.4

```sh
go mod download
go run .
```

### Con docker-compose

```sh
docker-compose -f "docker-compose.yml" up -d
```

## Uso

URL Local: http://localhost:5000
URL Local (Con Docker): http://localhost:5000 || http://conekta.amaurytq.localhost

URL desplegada: http://conekta.amaurytq.dev

## Tests

```sh
go test
```
