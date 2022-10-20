# FROM kneerunjun/gogingonic:latest
# Image to get golang app running on a light weight container 
# using alpine since its what makes a lightweight container 
FROM golang:1.18-alpine

ARG APPNAME
# git credentials are required only when the code is clone from the repo 
# we generally use Buddy and hence clone the code to the host's directory which is then volume mapped
# ARG gituname=$GIT_USERNAME
# ARG gitpassw=$GIT_PASSWORD

# RUN apk update && apk add git && apk add --no-cache openssh
RUN apk update 
RUN mkdir -p /usr/src/eensy /usr/bin/eensy /var/log/eensy 
WORKDIR /usr/src/eensy
COPY go.sum go.mod ./
RUN go mod download 
COPY . .
# here instead of relyingon the local machine for the code we are attempting to get the code from github
# RUN git clone https://gituname:gitpassw@github.com/eensymachines-in/usrv-accounts.git

RUN go build -o /usr/bin/eensy/${APPNAME} .