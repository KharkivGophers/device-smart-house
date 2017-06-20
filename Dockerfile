FROM golang
MAINTAINER Kharkiv Gophers (kostyamol@gmail.com)

COPY . /go/src/github.com/KharkivGophers/device-smart-house
WORKDIR /go/src/github.com/KharkivGophers/device-smart-house/cmd

RUN go get ./
RUN go build
CMD device-smart-house

#tcp conn for data from device to center
EXPOSE 3030

#tcp conn for config from device to center
EXPOSE 3000