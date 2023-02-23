FROM golang:1.19-alpine3.17 AS build

WORKDIR /app

COPY . .
RUN go mod download


RUN go build -o /service cmd/main.go
RUN touch /config.yml

## Deploy
FROM ubuntu:latest

WORKDIR /

RUN apt-get update
RUN apt-get upgrade
RUN apt-get install build-essential

COPY --from=build /service /service
COPY --from=build /config.yml /config.yml

#USER nonroot:nonroot
ENTRYPOINT ["/service", "-c", "/config.yml"]