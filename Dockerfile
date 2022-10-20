# FROM kneerunjun/gogingonic:latest
FROM golang:1.18-alpine

ARG APPNAME=accounts

RUN apk update && apk install git && echo apk --version
RUN mkdir -p /usr/src/eensy /usr/bin/eensy /var/log/eensy/accounts
WORKDIR /usr/src/eensy
# COPY go.sum go.mod ./
# RUN go mod download 
# COPY . .
# here instead of relyingon the local machine for the code we are attempting to get the code from github
RUN git clone git@github.com:eensymachines-in/usrv-accounts.git
RUN cd usrv-accounts
RUN git pull origin main
RUN go build -o /usr/bin/eensy/usrv-accounts .