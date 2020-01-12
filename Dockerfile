FROM golang:1.13.6-alpine3.11

WORKDIR /root
EXPOSE 9001

COPY src src

RUN cd src && go build main.go
RUN mv src/main .
RUN rm -rf src

CMD ./main