FROM golang:alpine AS builder

RUN mkdir -p /build
WORKDIR /build

COPY go.mod .
COPY . .
RUN go build -o main example1.go

RUN mkdir -p /dist
WORKDIR /dist
RUN cp /build/main ./main

EXPOSE 8001
CMD ["./main"]
