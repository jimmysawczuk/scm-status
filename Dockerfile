FROM golang:1.12 AS builder
WORKDIR /build
COPY . /go/src/github.com/jimmysawczuk/scm-status
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scm-status github.com/jimmysawczuk/scm-status

FROM alpine
RUN apk add --no-cache git tzdata
COPY --from=builder /build/scm-status /usr/bin/scm-status
WORKDIR /home
