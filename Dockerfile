FROM golang
MAINTAINER Kharkiv Gophers (kostyamol@gmail.com)

WORKDIR /home
COPY ./cmd/device-smart-house .

RUN \  
 chown daemon device-smart-house && \
 chmod +x device-smart-house
  
USER daemon
ENTRYPOINT ./device-smart-house
