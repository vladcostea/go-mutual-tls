FROM golang:1.13-alpine AS builder

RUN apk add --no-cache git

RUN mkdir service

ARG SERVICE

RUN mkdir -p service/.ssl
WORKDIR service

COPY go.mod .
COPY go.sum .

COPY cmd cmd
COPY secure secure

RUN go get ./...
RUN go build -ldflags '-w -s' -o ${SERVICE} cmd/http/${SERVICE}/main.go

FROM bash:latest

ARG SERVICE

COPY .ssl/ca.pem .ssl/service-${SERVICE}.pem .ssl/service-${SERVICE}-key.pem /.ssl/

COPY --from=builder /go/service/${SERVICE} .

ENV SERVICE ${SERVICE}
ENV CERT_DIR /.ssl

CMD ./${SERVICE}
