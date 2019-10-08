## Build
FROM golang:1.13.1-alpine3.10

WORKDIR /go/src/github.com/datawire/ingress-conformance/
COPY tools/echoserver.go echoserver.go
RUN go build -o echoserver .

## Run
FROM alpine:3.10

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/datawire/ingress-conformance/echoserver .

CMD ["./echoserver"]