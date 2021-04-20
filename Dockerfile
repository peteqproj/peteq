FROM golang:1.16-alpine3.12 AS dev

RUN apk update && apk add make gcc musl-dev

WORKDIR /peteq

COPY Makefile go.* /peteq/

RUN make dependency-update

COPY . .

RUN make test

RUN make compile

FROM alpine:3.12

RUN apk update && apk add ca-certificates

COPY --from=dev /peteq/dist /usr/local/bin

ENTRYPOINT [ "peteq-server" ]