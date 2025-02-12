FROM golang:1.23.6 as build

WORKDIR /app/

COPY go.mod go.sum /app/

RUN set -ex \
    && go mod download \
    && go mod verify

COPY . /app/

ENV CGO_ENABLED=0

RUN set -ex \
    && go build -trimpath -v -o ./bin/pushgateway-admin .


###############################################################################
FROM alpine:3

ENTRYPOINT ["/bin/pushgateway-admin"]

COPY --from=build --chown=nobody:nobody /app/bin/ /bin/

USER 65534
