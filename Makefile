REPO = gitlab.vrviu.com/clouddesk/cdp-admin-server

GIT_COMMIT := $(shell git show-branch --no-name HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "*" || true)
BUILD_VERSION := $(shell git describe --abbrev=10 --tags --always)
BUILD_TIME := $(shell date +%FT%T%z)

LDFLAGS := "\
-X \"${REPO}/common/version.buildGitCommit=${GIT_COMMIT} ${GIT_DIRTY}\" \
-X \"${REPO}/common/version.buildGitBranch=${GIT_BRANCH}\" \
-X \"${REPO}/common/version.buildVersion=${BUILD_VERSION}\" \
-X \"${REPO}/common/version.buildTime=${BUILD_TIME}\""

export CGO_ENABLED=0
# goctl
#GO_CTL_NAME=goctl1.4.1
GO_CTL_NAME=goctl

# go-zero生成代码风格
#GO_ZERO_STYLE=goZero #推荐使用这一种风格
GO_ZERO_STYLE=gozero


GO ?= go
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go")


all: release

.PHONY: test
test: # 运行项目测试
	go test -v --cover .//internal/..

.PHONY: fmt
fmt: # 格式化代码
	$(GOFMT) -w $(GOFILES)

.PHONY: api
api: # 生成 cdp-admin-service 的代码
	$(GO_CTL_NAME) api go -api cdp-admin-service.api -dir .  --style=$(GO_ZERO_STYLE)
	@echo "Generate cdp-admin-service files successfully"


swagger:
	goctl api plugin -plugin goctl-swagger="swagger -filename cdp-admin-service.json" -api cdp-admin-service.api -dir .

.PHONY: mysql
mysql: # 生成 cdp-admin-service mysql 模块代码
	$(GO_CTL_NAME) model mysql ddl  --cache=false -d /model/ -s doc/sql/*.sql  --style=$(GO_ZERO_STYLE)
	@echo "Generate cdp-admin-service files successfully"


.PHONY: all
all: # 生成全部api和rpc代码
	#make api
	#make gen-admin-rpc
	@echo "Generate all files successfully"

.PHONY: release
release: # 编译生成可执行文件
	go mod tidy 
	go vet ./...
	GOWORK=off go build -ldflags $(LDFLAGS)  cdp-admin-service.go 


.PHONY: help
help: # 显示帮助
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

.PHONY: precommit
precommit: # 预提交
	make api
	make swagger
	go mod tidy 
	go fmt ./...
	GOWORK=off go build -ldflags $(LDFLAGS)  cdp-admin-service.go