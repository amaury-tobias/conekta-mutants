FROM golang:alpine as build
WORKDIR /src
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o /bin/main .

FROM scratch
COPY --from=build /bin/main /bin/main
EXPOSE 5000
ENTRYPOINT ["/bin/main"]