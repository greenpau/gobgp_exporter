.PHONY: test clean qtest
APP_VERSION:=1.0.0
GIT_COMMIT:=$(shell git describe --dirty --always)
GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD -- | head -1)
BUILD_USER:=$(shell whoami)
BUILD_DATE:=$(shell date +"%Y-%m-%d")
BINARY:=gobgp_exporter
VERBOSE:=-v

all:
	@echo "Version: $(APP_VERSION), Branch: $(GIT_BRANCH), Revision: $(GIT_COMMIT)"
	@echo "Build on $(BUILD_DATE) by $(BUILD_USER)"
	@cat gobgp.pb.extended.patch > $(GOPATH)/src/github.com/osrg/gobgp/api/gobgp.pb.extended.go
	@CGO_ENABLED=0 go build -o ./$(BINARY) $(VERBOSE) \
		-ldflags="-w -s \
		-X github.com/prometheus/common/version.Version=$(APP_VERSION) \
		-X github.com/prometheus/common/version.Revision=$(GIT_COMMIT) \
		-X github.com/prometheus/common/version.Branch=$(GIT_BRANCH) \
		-X github.com/prometheus/common/version.BuildUser=$(BUILD_USER) \
		-X github.com/prometheus/common/version.BuildDate=$(BUILD_DATE) \
		-X main.appVersion=$(APP_VERSION) \
		-X main.gitBranch=$(GIT_BRANCH) \
		-X main.gitCommit=$(GIT_COMMIT) \
		-X main.buildUser=$(BUILD_USER) \
		-X main.buildDate=$(BUILD_DATE)" \
		-gcflags="all=-trimpath=$(GOPATH)/src" \
		-asmflags="all=-trimpath $(GOPATH)/src"
	@echo "Done!"

test: all
	@go test -v ./*.go
	@echo "PASS: core tests"
	@echo "OK: all tests passed!"

clean:
	@rm -rf bin/
	@echo "OK: clean up completed"

qtest:
	@./$(BINARY) -version
	@./$(BINARY) -web.listen-address 0.0.0.0:5000

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
