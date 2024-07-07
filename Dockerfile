FROM golang:1.22 AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ./cmd/main ./cmd/main.go

#Deploy
FROM ubuntu:22.04

WORKDIR /app

COPY --from=build /app/cmd/main ./cmd/main 
COPY --from=build /app/configs/ ./configs/
COPY --from=build /app/migrations/ ./migrations/

EXPOSE 8080

ENV MODE=DOCKER

ENTRYPOINT ["./cmd/main"]

