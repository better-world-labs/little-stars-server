FROM golang:1.19-alpine3.16 as builder
MAINTAINER shenweijie <shenweijie@openviewtech.com>

RUN sed 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' -i /etc/apk/repositories && \
apk update \
    && apk add gcc

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

#依赖准备
WORKDIR /src
COPY go.mod .
RUN --mount=type=cache,target=/root/.cache/go-download go mod download


COPY assert assert
COPY cmd cmd
COPY internal internal

RUN go mod tidy

#生成
RUN go run -tags dynamic ./cmd/tools/autoload -s internal/controller -s internal/service -p command -f autoLoad -o internal/command/autoload.go

#编译
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags "-s -w" -tags musl -o /src/bin/app /src/cmd/app


FROM alpine:3.16
ARG ENVIRONMENT=dev
ENV env=${ENVIRONMENT}

RUN sed 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' -i /etc/apk/repositories && \
apk update \
    && apk add tzdata\
    && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime -f

WORKDIR /app

COPY assert assert
COPY config config

COPY --from=builder /src/bin/app /app
CMD ./app -e $env

EXPOSE 8080