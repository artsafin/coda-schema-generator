FROM golang:1.17-alpine as deps

ADD go.mod /app/go.mod
WORKDIR /app

RUN go mod download
RUN apk add alpine-sdk




FROM deps as build

ADD . /app
WORKDIR /app

RUN go mod tidy && \
    time go build -o "/tmp/csg" ./cmd && \
    chmod a+x /tmp/csg && \
    echo -n "BIN SIZE: " && du -k /tmp/csg




FROM alpine

RUN addgroup -g 9999 -S user && \
    adduser -u 9999 -G user -S -H user

COPY --from=build /tmp/csg /
ENTRYPOINT ["/csg"]

USER user
