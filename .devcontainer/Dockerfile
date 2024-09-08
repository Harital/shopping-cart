# syntax=docker/dockerfile:1

########################################
### BASE IMAGE #########################
########################################

# Choosing golang alpine due to small footprint
ARG GO_VERSION
ARG ALPINE_VERSION

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
RUN gofmt -w -s .


########################################
### DEVELOPMENT ########################
########################################

ARG ALPINE_VERSION

FROM alpine:${ALPINE_VERSION} AS development

WORKDIR $GOPATH/src

# Declare variables to be used in docker
ARG USERNAME
ENV PATH=/usr/local/go/bin:/usr/local/go/:/home/${USERNAME}:/go/bin:${PATH}

# Copy go packages and compiled code
COPY --from=base /usr/local/go/ /usr/local/go/

# Update the alpine image packages to avoid possible Anchore issues and
# add developer user and add it to sudoers file
RUN apk -U upgrade && \
    apk add --no-cache sudo&& \
    addgroup -S ${USERNAME} && \
    adduser -S -G ${USERNAME}  ${USERNAME} && \
    mkdir -p /etc/sudoers.d && \
    echo ${USERNAME} ALL=\(ALL\) NOPASSWD:ALL > /etc/sudoers.d/${USERNAME} && \
    chmod 0440 /etc/sudoers.d/${USERNAME}

# Avoid VSCode extension reinstall
RUN mkdir -p /home/${USERNAME}/.vscode-server/extensions \
    /home/${USERNAME}/.vscode-server-insiders/extensions && \
    chown -R ${USERNAME} \
    /home/${USERNAME}/.vscode-server \
    /home/${USERNAME}/.vscode-server-insiders

# Install required development packages, not available in original alpine distro
# permissions in go foder should be open enough to let developer user read and write
ENV GOPATH=/go
RUN apk add --no-cache bash git openssh libc-dev gcc && \
    umask 0000 && \
    /usr/local/go/bin/go install -v golang.org/x/tools/gopls@v0.16.1 && \
    /usr/local/go/bin/go install -v github.com/go-delve/delve/cmd/dlv@v1.23.0 && \
    /usr/local/go/bin/go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1 && \
    /usr/local/go/bin/go install -v golang.org/x/tools/cmd/godoc@v0.24.0 && \
    /usr/local/go/bin/go install -v github.com/cweill/gotests/gotests@v1.6.0 && \
    /usr/local/go/bin/go install -v github.com/air-verse/air@v1.52.3 && \
    /usr/local/go/bin/go install -v go.uber.org/mock/mockgen@v0.4.0

EXPOSE 8080

CMD ["/bin/bash"]

