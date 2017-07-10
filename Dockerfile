FROM alpine
MAINTAINER Kharkiv Gophers (kostyamol@gmail.com)

RUN mkdir -p /home/device-smart-house/bin

WORKDIR /home/device-smart-house/bin
COPY ./cmd/device-smart-house .

RUN \  
    chown daemon device-smart-house && \
    chmod +x device-smart-house
    
USER daemon
ENTRYPOINT ["./device-smart-house"]
CMD ["fridge", "LG", "FF:FF:FF:FF:FF:FF"]
