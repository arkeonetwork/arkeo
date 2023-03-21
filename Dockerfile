#
# Arkeo
#

ARG GO_VERSION="1.20"

#
# Build
#
FROM golang:${GO_VERSION} as builder

ARG GIT_VERSION
ARG GIT_COMMIT

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux

# Download go dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TAG=testnet
RUN make install

#
# Main
#
FROM ubuntu:kinetic

RUN apt-get update -y && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
      jq=1.6-2.1ubuntu3 curl=7.85.0-1ubuntu0.3 htop=3.2.1-1 vim=2:9.0.0242-1ubuntu1 ca-certificates=20211016ubuntu0.22.10.1 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN update-ca-certificates

# Copy the compiled binaries over.
COPY --from=builder /go/bin/sentinel /go/bin/arkeod /usr/bin/

COPY scripts /scripts

ARG TAG=testnet
ENV NET=$TAG

# default to fullnode
CMD ["arkeod", "start"]
