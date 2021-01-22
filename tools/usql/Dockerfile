FROM golang:alpine as dev

RUN apk update && \
    apk add git curl gcc musl-dev

RUN git clone https://github.com/xo/usql && \
    cd usql && go build -tags 'postgres' && \
    chmod +x usql

FROM alpine

COPY --from=dev /go/usql/usql /usr/local/bin/usql

CMD [ "usql" ]
