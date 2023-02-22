#
# Arkeo
#

ARG GO_VERSION="1.18"

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

RUN apt-get update -y && apt-get upgrade -y && apt-get install -y jq curl htop vim

# Copy the compiled binaries over.
COPY --from=builder /go/bin/sentinel /go/bin/arkeod /usr/bin/

COPY scripts /scripts

ARG TAG=testnet
ENV NET=$TAG

# default to fullnode
ENTRYPOINT ["arkeod", "start", "--home", "/.arkeo"]
