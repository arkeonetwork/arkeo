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

RUN curl https://get.ignite.com/cli! | bash

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    curl git jq vim make protobuf-compiler xz-utils sudo python3-pip \
    && rm -rf /var/cache/apt/lists \
    && go install mvdan.cc/gofumpt@v0.3.0

# Download go dependencies
WORKDIR /app
COPY go.mod go.sum ./

COPY . .

ARG TAG=mainnet
RUN ignite chain build --proto-all-modules
# RUN ./scripts/protocgen.sh
RUN make install
RUN ls 
RUN ls /go/bin
RUN ls /usr/bin

#
# Main
#
FROM golang:${GO_VERSION}-alpine

# Copy the compiled binaries over.
COPY --from=builder /go/bin/sentinel /go/bin/arkeod /usr/bin/

COPY scripts /scripts

# default to mainnet
ARG TAG=mainnet
ENV NET=$TAG

# default to fullnode
CMD ["/scripts/run-arkeo.sh"]
