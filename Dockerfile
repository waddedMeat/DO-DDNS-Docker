FROM golang:1.16-alpine as builder
ARG TARGETPLATFORM

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/ddodns/mupdater/
COPY . .
RUN go mod download
RUN go mod verify


RUN export CGO_ENABLED=0; \
    export GOOS=linux; \
    [ "$TARGETPLATFORM" == "linux/amd64" ] && export GOARCH=amd64; \
    [ "$TARGETPLATFORM" == "linux/arm64" ] && export GOARCH=arm64; \
    [ "$TARGETPLATFORM" == "linux/arm/v7" ] && export GOARCH=arm; \
    [ -z $GOARCH ] && export GOARCH=amd64; \
    go build -ldflags="-w -s" -o /go/bin/updater;

##
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/updater /go/bin/updater

ENTRYPOINT ["/go/bin/updater"]