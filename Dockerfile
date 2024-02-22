##### BUILDER #####

## use golang 1.18 alpine image as base from https://hub.docker.com/_/golang/tags to build the project
FROM golang:1.18-alpine3.16 as builder

## copy source files from the project in the image
COPY . /src
WORKDIR /src

## clean and build the go project
RUN go mod tidy -go=1.18

ENV GOOS="linux"
ENV GOARCH="amd64"
ENV CGO_ENABLED="0"

RUN go build -ldflags="-s -w" -o service-shop-go cmd/main.go

## create a new smaller image with alpine:3.16 as base and copy the built project from the previous image
FROM alpine:3.16

ENV IMG_VERSION="1.0.0"

COPY --from=builder /src/service-shop-go /
COPY --from=builder /src/config/service.yaml /config/

ENTRYPOINT ["/service-shop-go"]
CMD ["--config","/config/service.yaml"]

EXPOSE 8080 8443

HEALTHCHECK --interval=30s --timeout=5s --retries=3 --start-period=10s \
    CMD wget -q -T 5 --spider http://localhost:8080/health/healthiness

## image description
LABEL org.opencontainers.image.title="Shop Service" \
    org.opencontainers.image.description="Azubi shop service" \
    org.opencontainers.image.version="${IMG_VERSION}" \
    org.opencontainers.image.source="https://bitbucket.easy.de/users/n.gauche/repos/service-shop-go.git" \
    org.opencontainers.image.vendor="EASY SOFTWARE AG (www.easy-software.com)" \
    org.opencontainers.image.authors="EASY SOFTWARE AG" \
    maintainer="EASY SOFTWARE AG" \
    NAME="shop-service"