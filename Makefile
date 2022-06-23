help:
	@echo "make gen 生成internal/command/autoload.go"
	@echo "make dev 监听文件修改，动态生成internal/command/autoload.go"
	@echo "make build 编译项目"
gen:
	go run -tags dynamic ./cmd/tools/autoload -s internal/controller -s internal/service -p command -f autoLoad -o internal/command/autoload.go
dev:
	go run -tags dynamic ./cmd/tools/autoload -w -s internal/controller -s internal/service -p command -f autoLoad -o internal/command/autoload.go

build:
	make gen
	go build -tags dynamic -o dist/server aed-api-server/cmd/app