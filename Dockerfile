FROM golang:1.13.6-alpine3.11

WORKDIR /root
EXPOSE 9001

COPY src src

RUN go build -o patch-hosts-service-linux-amd64 src/main/main.go

CMD ./patch-hosts-service-linux-amd64