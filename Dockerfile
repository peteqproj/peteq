FROM golang:1.15-alpine3.12 AS dev

RUN apk add make

WORKDIR /peteq

COPY Makefile go.* /peteq/

RUN make dependency-update

COPY . .

RUN make build

FROM alpine:3.12

RUN apk update && apk add ca-certificates

COPY --from=dev /peteq/dist/peteq /peteq

CMD [ "/peteq" ]