FROM golang:1.19-alpine3.17 AS build

RUN apk add build-base
WORKDIR /app

COPY . .
RUN go mod download


RUN go build -o /service cmd/main.go
RUN touch /config.yml

## Deploy
FROM alpine:latest

WORKDIR /

COPY --from=build /service /service
COPY --from=build /config.yml /config.yml

#USER nonroot:nonroot
ENTRYPOINT ["/service", "-c", "/config.yml"]