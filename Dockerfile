FROM golang:alpine as builder
RUN mkdir /build
RUN apk add bash make
ADD . /build/
WORKDIR /build
RUN make

FROM alpine
RUN mkdir -p /opt
COPY --from=builder /build/ingress-controller-conformance /opt/
WORKDIR /opt
ENTRYPOINT ["/opt/ingress-controller-conformance"]

