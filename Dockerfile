FROM gitlab.openviewtech.com:5050/devops/alpine-golang as builder
MAINTAINER shenweijie <shenweijie@openviewtech.com>

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

#依赖准备
WORKDIR /src
COPY go.mod .
RUN --mount=type=cache,target=/root/.cache/go-download go mod download

#编译
COPY . .
RUN go mod tidy
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags "-s -w" -tags musl -o /src/bin/app /src/cmd/app


FROM gitlab.openviewtech.com:5050/devops/alpine-timezone-shanghai
ARG ENVIRONMENT=testing

WORKDIR /app
RUN ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime -f

COPY assert assert
COPY config-${ENVIRONMENT}.yaml config.yaml

COPY --from=builder /src/bin/app /app
CMD ./app run -c config.yaml

EXPOSE 8080