# syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go mod download

RUN go build -o /fakebee

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /fakebee /fakebee

ENTRYPOINT [ "/fakebee" ]