FROM golang:1.19.0-alpine3.16 as builder

RUN apk add --no-cache git alpine-sdk

ADD . /go/src/admincheckapi
WORKDIR /go/src/admincheckapi

# build the source
RUN make tidy && make build

# use a minimal alpine image
FROM alpine:3.16

# add ca-certificates in case you need them
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# set working directory
WORKDIR /go/bin

COPY --from=builder /go/src/admincheckapi .
COPY --from=builder /go/src/admincheckapi/config.yaml .

USER 1001
EXPOSE 1234/tcp

# run the binary
CMD ["./admincheckapi"]
