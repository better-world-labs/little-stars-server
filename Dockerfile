FROM gitlab.openviewtech.com:5050/devops/alpine-golang as builder
MAINTAINER shenweijie <shenweijie@openviewtech.com>

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

#依赖准备
WORKDIR /src
COPY go.mod .
RUN --mount=type=cache,target=/root/.cache/go-download go mod download


COPY assert assert
COPY cmd cmd
COPY internal internal

RUN go mod tidy -compat=1.17

#生成
RUN go run -tags dynamic ./cmd/tools/autoload -s internal/controller -s internal/service -p command -f autoLoad -o internal/command/autoload.go

#编译
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags "-s -w" -tags musl -o /src/bin/app /src/cmd/app


FROM gitlab.openviewtech.com:5050/devops/alpine-timezone-shanghai
ARG ENVIRONMENT=dev
ENV env=${ENVIRONMENT}

WORKDIR /app
RUN ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime -f

COPY assert assert
COPY config config

COPY --from=builder /src/bin/app /app
CMD ./app -e $env

EXPOSE 8080