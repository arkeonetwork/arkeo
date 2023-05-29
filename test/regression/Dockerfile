FROM golang:1.20

# hadolint ignore=DL3008,DL4006
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    curl git jq make protobuf-compiler xz-utils sudo postgresql postgresql-contrib python3-pip \
    && rm -rf /var/cache/apt/lists \
    && go install mvdan.cc/gofumpt@v0.3.0 \
    && go install github.com/jackc/tern/v2@latest

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN mkdir /regtest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build make _build-test-regression
