FROM golang:1.12 AS builder
ENV GOPATH=/gocode
RUN mkdir -p /test/uapp
WORKDIR /test/uapp
COPY . .
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOPROXY=https://gocenter.io
RUN go build -o /test/uapp/uapp .

FROM alpine

COPY --from=builder /test/uapp /test/uapp

RUN apk update && apk add ca-certificates

CMD ["/test/uapp/uapp"]