# syntax=docker/dockerfile:1

########################################
### BASE IMAGE #########################
########################################

# Choosing golang alpine due to small footprint
ARG GO_VERSION=1.22
ARG ALPINE_VERSION=3.20

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base

ENV CGO_ENABLED=0 \
    GIN_MODE=release

# Update the golang-alpine image packages to avoid possible Anchore issues.
RUN apk -U upgrade

# Go packages setup
WORKDIR $GOPATH/src
COPY go.* ./
RUN go install ./...

# Build
COPY . ./
RUN gofmt -w -s . && \
    CGO_ENABLED=${CGO_ENABLED} GIN_MODE=${GIN_MODE} go build -o shopping-cart cmd/main.go


########################################
### PRODUCTION #########################
########################################

ARG ALPINE_VERSION

FROM alpine:${ALPINE_VERSION} AS production

# Declare variables to be used in docker
ARG USER_UID=1002
ARG USER_GID=${USER_UID}

# Update the alpine image packages to avoid possible Anchore issues.
RUN apk -U upgrade

RUN apk add --no-cache bash git wget nodejs npm openssh libc-dev gcc

# Download ajv-cli
RUN npm install -g ajv-cli


COPY --from=base /go/src/shopping-cart .

# Required. Run service as non-root.
USER ${USER_UID}


CMD ["./shopping-cart"]
