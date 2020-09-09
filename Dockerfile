# ENV VARS:
# 
# - ERIKBOTDEV_CONFIG_FILE_NAME
# - TWITCH_CLIENT_ID
# - TWITCH_CLIENT_SECRET
# - TWITCH_OAUTH_TOKEN
# - PORT
#
# You can override the Go version used to build the image.
# See project Makefile if using make.
# See docker --build-arg if building directly.
ARG GOLANG_VERSION=1.15-buster
ARG ALPINE_VERSION=3.11.5

FROM golang:${GOLANG_VERSION} AS builder

RUN apt update && apt install -y nodejs npm

ENV GOPATH="/go"
WORKDIR $GOPATH/src/github.com/erikstmartin/erikbotdev

COPY . .

RUN cd web && npm run build

RUN GO111MODULE=on CGO_ENABLED=0 GOPROXY="https://proxy.golang.org" go build -o erikbotserver ./cmd/server

FROM alpine:${ALPINE_VERSION}

ENV GO111MODULE=on

RUN mkdir -p $HOME/erikbotserver
WORKDIR $HOME/erikbotserver

COPY --from=builder /go/src/github.com/erikstmartin/erikbotdev .

# Add tini, see https://github.com/gomods/athens/issues/1155 for details.
EXPOSE 3000

CMD ["./erikbotserver", "serve"]
