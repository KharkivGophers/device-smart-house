FROM golang
MAINTAINER Kharkiv Gophers (kostyamol@gmail.com)

COPY . /go/src/github.com/KharkivGophers/device-smart-house
WORKDIR /go/src/github.com/KharkivGophers/device-smart-house/cmd

RUN go get ./
RUN go build -o device-smart-house
CMD device-smart-house
