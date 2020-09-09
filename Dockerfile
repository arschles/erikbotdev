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

WORKDIR $GOPATH/src/github.com/erikstmartin/erikbotdev

COPY . .
RUN mv examples /bin/configs

ARG VERSION="unset"

RUN cd web && npm run build && mkdir -p /bin/public/build && mv public/build /bin/public/build

RUN GO111MODULE=on CGO_ENABLED=0 GOPROXY="https://proxy.golang.org" go build -o /bin/erikbotserver ./cmd/server

FROM alpine:${ALPINE_VERSION}

ENV GO111MODULE=on

RUN mkdir -p $HOME/erikbotserver/web/public/build && mkdir -p $HOME/erikbotserver/configs && mkdir /configs
WORKDIR $HOME/erikbotserver

COPY --from=builder /bin/public/build /web/public/build 
COPY --from=builder /bin/erikbotserver ./erikbotserver
COPY --from=builder /bin/configs /configs
RUN ls /configs

# Add tini, see https://github.com/gomods/athens/issues/1155 for details.
EXPOSE 3000

CMD ["./erikbotserver"]
