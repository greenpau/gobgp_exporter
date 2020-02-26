.PHONY: test clean qtest deploy dist linter dep gobgp-down gobgp tag
APP_VERSION:=$(shell cat VERSION | head -1)
GIT_COMMIT:=$(shell git describe --dirty --always)
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD -- | head -1)
BUILD_USER:=$(shell whoami)
BUILD_DATE:=$(shell date +"%Y-%m-%d")
ALIAS=gobgp_exporter
BINARY:=gobgp-exporter
VERBOSE:=-v
PROJECT=github.com/ovnworks/$(ALIAS)
PKG_DIR=pkg/$(ALIAS)
PKG_PORT="9474"

all:
	@echo "Version: $(APP_VERSION), Branch: $(GIT_BRANCH), Revision: $(GIT_COMMIT)"
	@echo "Build on $(BUILD_DATE) by $(BUILD_USER)"
	@mkdir -p bin/
	@rm -rf ./bin/*
	@CGO_ENABLED=0 go build -o ./bin/$(BINARY) $(VERBOSE) \
		-gcflags="all=-trimpath=$(GOPATH)/src" \
		-asmflags="all=-trimpath $(GOPATH)/src" \
		-ldflags="-w -s -v \
		-X github.com/prometheus/common/version.Version=$(APP_VERSION) \
		-X github.com/prometheus/common/version.Revision=$(GIT_COMMIT) \
		-X github.com/prometheus/common/version.Branch=$(GIT_BRANCH) \
		-X github.com/prometheus/common/version.BuildUser=$(BUILD_USER) \
		-X github.com/prometheus/common/version.BuildDate=$(BUILD_DATE) \
		-X $(PROJECT)/$(PKG_DIR).appName=$(BINARY) \
		-X $(PROJECT)/$(PKG_DIR).appVersion=$(APP_VERSION) \
		-X $(PROJECT)/$(PKG_DIR).gitBranch=$(GIT_BRANCH) \
		-X $(PROJECT)/$(PKG_DIR).gitCommit=$(GIT_COMMIT) \
		-X $(PROJECT)/$(PKG_DIR).buildUser=$(BUILD_USER) \
		-X $(PROJECT)/$(PKG_DIR).buildDate=$(BUILD_DATE)" \
		./cmd/$(ALIAS)/*.go
	@echo "Done!"

linter:
	@#go get -u golang.org/x/lint/golint
	@golint ./$(PKG_DIR)/*.go
	@echo "PASS: golint"

test: linter all
	@./bin/$(BINARY) -metrics
	@go test -v ./$(PKG_DIR)/*.go
	@echo "PASS: core tests"
	@echo "OK: all tests passed!"

clean:
	@rm -rf bin/
	@rm -rf dist/
	@echo "OK: clean up completed"

deploy:
	@sudo rm -rf /usr/bin/$(BINARY)
	@sudo cp ./bin/$(BINARY) /usr/bin/$(BINARY)

qtest:
	@./bin/$(BINARY) -version
	@./bin/$(BINARY) \
		-web.listen-address 0.0.0.0:$(PKG_PORT) \
		-log.level debug \
		-gobgp.poll-interval 5

dist: all
	@mkdir -p ./dist
	@rm -rf ./dist/*
	@mkdir -p ./dist/$(BINARY)-$(APP_VERSION).linux-amd64
	@cp ./bin/$(BINARY) ./dist/$(BINARY)-$(APP_VERSION).linux-amd64/
	@cp ./README.md ./dist/$(BINARY)-$(APP_VERSION).linux-amd64/
	@cp LICENSE ./dist/$(BINARY)-$(APP_VERSION).linux-amd64/
	@cp assets/systemd/install.sh ./dist/$(BINARY)-$(APP_VERSION).linux-amd64/install.sh
	@cp assets/systemd/uninstall.sh ./dist/$(BINARY)-$(APP_VERSION).linux-amd64/uninstall.sh
	@chmod +x ./dist/$(BINARY)-$(APP_VERSION).linux-amd64/*.sh
	@cd ./dist/ && tar -cvzf ./$(BINARY)-$(APP_VERSION).linux-amd64.tar.gz ./$(BINARY)-$(APP_VERSION).linux-amd64

dep:
	@echo "Making dependencies check ..."
	@#echo "Clean GOPATH/pkg/dep/sources/ if necessary"
	@#rm -rf $GOPATH/pkg/dep/sources/https---github.com-ovnworks*
	@dep version || go get -u github.com/golang/dep/cmd/dep
	@dep ensure

gobgp:
	@sudo systemctl stop gobgpd
	@sudo systemctl start gobgpd
	@sleep 1
	@sudo systemctl status gobgpd
	@sudo gobgp global rib add -a ipv4 10.10.10.0/24 origin igp
	@sudo gobgp global rib add -a ipv4 10.10.20.0/24 origin igp
	@sudo gobgp global rib add -a ipv4 10.10.30.0/24 origin igp
	@sudo gobgp global rib

gobgp-down:
	@sudo systemctl stop gobgpd

tag:
	@git tag -s "v$(APP_VERSION)" -m "v$(APP_VERSION)"
	@git push --tags
