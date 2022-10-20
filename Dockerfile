# FROM kneerunjun/gogingonic:latest
FROM golang:1.18-alpine

ARG APPNAME=accounts

RUN mkdir -p /usr/src/eensy/accounts /usr/bin/eensy /var/log/eensy/accounts
WORKDIR /usr/src/eensy/accounts
COPY go.sum go.mod ./
RUN go mod download 
COPY . .
RUN go build -o /usr/bin/eensy/accounts .