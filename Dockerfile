# Build image
FROM golang:1.14-alpine AS builder

RUN apk update \
    && apk upgrade \
    && apk add --update \
    ca-certificates \
    gcc \
    git \
    libc-dev \
    make \
    && update-ca-certificates

WORKDIR $GOPATH/src/github.com/guilherme-santos/messageboard

COPY go.mod go.sum ./
RUN go mod download

ENV GIT_TAG $GIT_TAG
ENV GIT_COMMIT $GIT_COMMIT

COPY . ./
RUN make go-install

# Final image
FROM alpine:3.11

LABEL maintainer="Guilherme Silveira <xguiga@gmail.com>"

# set up nsswitch.conf for Go's "netgo" implementation
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf

ENV HEALTHCHECK_VERSION 1.0.0
ENV HEALTHCHECK_URL https://github.com/gioxtech/healthcheck/releases/download/v${HEALTHCHECK_VERSION}/healthcheck-${HEALTHCHECK_VERSION}
RUN wget ${HEALTHCHECK_URL} -O /usr/bin/healthcheck && \
    chmod +x /usr/bin/healthcheck

HEALTHCHECK --start-period=5s --interval=30s --timeout=1s --retries=3 CMD healthcheck -http-addr http://localhost/ping

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/messageboard /usr/bin/
COPY --from=builder /go/src/github.com/guilherme-santos/messageboard/messages.csv /etc/messageboard/

EXPOSE 80

CMD ["messageboard"]
