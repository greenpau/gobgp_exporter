APP_VERSION:=$(shell cat VERSION | head -1)
GIT_COMMIT:=$(shell git describe --dirty --always)
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD -- | head -1)
BUILD_USER:=$(shell whoami)
BUILD_DATE:=$(shell date +"%Y-%m-%d")
ALIAS=gobgp_exporter
BINARY:=gobgp-exporter
VERBOSE:=-v
PROJECT=github.com/greenpau/$(ALIAS)
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
	@echo "DEBUG: completed $@"

.PHONY: linter
linter:
	@echo "Running lint checks"
	@golint ./...
	@echo "DEBUG: completed $@"

.PHONY: covdir
covdir:
	@echo "Creating .coverage/ directory"
	@mkdir -p .coverage
	@echo "DEBUG: completed $@"

.PHONY: test
test: covdir linter
	@./bin/$(BINARY) -metrics
	@go test $(VERBOSE) -coverprofile=.coverage/coverage.out ./...
	@echo "DEBUG: completed $@"

.PHONY: clean
clean:
	@rm -rf bin/
	@rm -rf dist/
	@echo "DEBUG: completed $@"

.PHONY: deploy
deploy:
	@sudo rm -rf /usr/bin/$(BINARY)
	@sudo cp ./bin/$(BINARY) /usr/bin/$(BINARY)
	@echo "DEBUG: completed $@"

.PHONY: qtest
qtest:
	@./bin/$(BINARY) -version
	@./bin/$(BINARY) \
		-web.listen-address 0.0.0.0:$(PKG_PORT) \
		-log.level debug \
		-gobgp.poll-interval 5
	@echo "DEBUG: completed $@"

.PHONY: dist
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
	@echo "DEBUG: completed $@"

.PHONY: dep
dep:
	@echo "Making dependencies check ..."
	@versioned || go install github.com/greenpau/versioned/cmd/versioned@v1.0.28
	@golint || go install golang.org/x/lint/golint@latest
	@echo "DEBUG: completed $@"

.PHONY: gobgp
gobgp:
	@sudo systemctl stop gobgpd
	@sudo systemctl start gobgpd
	@sleep 1
	@sudo systemctl status gobgpd
	@sudo gobgp global rib add -a ipv4 10.10.10.0/24 origin igp
	@sudo gobgp global rib add -a ipv4 10.10.20.0/24 origin igp
	@sudo gobgp global rib add -a ipv4 10.10.30.0/24 origin igp
	@sudo gobgp global rib
	@echo "DEBUG: completed $@"

.PHONY: gobgp-down
gobgp-down:
	@sudo systemctl stop gobgpd
	@echo "DEBUG: completed $@"

.PHONY: release
release:
	@echo "Making release"
	@go mod tidy
	@go mod verify
	@if [ $(GIT_BRANCH) != "main" ]; then echo "cannot release to non-main branch $(GIT_BRANCH)" && false; fi
	@git diff-index --quiet HEAD -- || ( echo "git directory is dirty, commit changes first" && false )
	@versioned -patch
	@echo "Patched version"
	@git add VERSION
	@git commit -m "released v`cat VERSION | head -1`"
	@git tag -a v`cat VERSION | head -1` -m "v`cat VERSION | head -1`"
	@git push
	@git push --tags
	@echo "DEBUG: completed $@"