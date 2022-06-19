# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go mod download

RUN go build -o /fakebee

ENTRYPOINT [ "/fakebee" ]