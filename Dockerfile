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

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN curl https://get.ignite.com/cli! | bash

# nolint
# RUN apt-get update \
    # && apt-get install -y --no-install-recommends \
    # curl git jq vim make protobuf-compiler xz-utils sudo python3-pip \
    # && rm -rf /var/cache/apt/lists \
    # && go install mvdan.cc/gofumpt@v0.3.0

# Download go dependencies
WORKDIR /app
COPY go.mod go.sum ./

COPY . .

RUN ignite chain build --proto-all-modules
ARG TAG=mainnet
RUN make install

#
# Main
#
FROM golang:${GO_VERSION}-alpine

# Copy the compiled binaries over.
COPY --from=builder /go/bin/sentinel /go/bin/arkeod /usr/bin/

COPY scripts /scripts

ENTRYPOINT ["/scripts/genesis.sh"]

ARG TAG=mainnet
ENV NET=$TAG

# default to fullnode
CMD ["/scripts/run-arkeo.sh"]
