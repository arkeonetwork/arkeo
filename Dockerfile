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
ARG TAG=mainnet
RUN make install

#
# Main
#
FROM golang:${GO_VERSION}-alpine

RUN apk add --no-cache \
    jq=1.6-r2 \
    curl=7.87.0-r1

# Copy the compiled binaries over.
COPY --from=builder /go/bin/sentinel /go/bin/arkeod /usr/bin/

COPY scripts /scripts

ENTRYPOINT ["/scripts/genesis.sh"]

ARG TAG=mainnet
ENV NET=$TAG

# default to fullnode
CMD ["/scripts/run-arkeo.sh"]
