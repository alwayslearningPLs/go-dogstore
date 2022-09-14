# golang:alpine with version 3.16 (last alpine version also)
FROM golang@sha256:c9a90742f5457fae80d8f9f1c9fc6acd6884c749dc6c5b11c22976973564dd4f as base

RUN mkdir --parent /go/src/github.com/MrTimeout/spacetrack && \
  apk add --no-cache --update ca-certificates tzdata git && update-ca-certificates

WORKDIR /go/src/github.com/MrTimeout/rest-api

COPY . .
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it.sh

RUN go mod download && go mod verify && \
  CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-w -s" -o /go/bin ./... && \
  chmod +x /wait-for-it.sh

FROM alpine:3.16

RUN apk add --no-cache bash

COPY --from=base /go/bin/wso2-rest-api-example /wso2-rest-api-example
COPY --from=base /wait-for-it.sh /

CMD ["/wso2-rest-api-example"]
