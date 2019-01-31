# BUILDER

FROM golang:1.11 AS builder
RUN go version

COPY  . /go/src/github.com/gravitational/slackbot
WORKDIR /go/src/github.com/gravitational/slackbot

RUN set -x && \
    go get github.com/golang/dep/cmd/dep && \
    dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o slackbot .

# RUNTIME
FROM quay.io/gravitational/debian-tall

WORKDIR /bot/

COPY --from=builder /go/src/github.com/gravitational/slackbot/slackbot ./

ENTRYPOINT [ "./slackbot" ]
