FROM golang:1.15-alpine3.13 as builder
WORKDIR /src
COPY . /src/

RUN go build -o ga-freshrelease


FROM alpine:3.13

COPY --from=builder /src/ga-freshrelease /usr/bin/ga-freshrelease

ENTRYPOINT ["/usr/bin/ga-freshrelease"]
