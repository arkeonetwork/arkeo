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
RUN make install && go install github.com/jackc/tern@latest

#
# Main
#
FROM ubuntu:kinetic

# hadolint ignore=DL3008,DL4006
RUN apt-get update -y && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
      jq curl htop vim ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN update-ca-certificates

# Copy the compiled binaries over.
COPY --from=builder /go/bin/sentinel /go/bin/arkeod /go/bin/indexer /go/bin/api /go/bin/tern /usr/bin/
COPY scripts /scripts

ARG TAG=testnet
ENV NET=$TAG

ENTRYPOINT ["scripts/genesis.sh"]

# default to fullnode
CMD ["arkeod", "start"]
