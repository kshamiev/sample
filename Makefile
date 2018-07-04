APP=sample

DIR=$(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

GOPATH := $(DIR):$(GOPATH)
DATE=$(shell date -u +%Y%m%d.%H%M%S.%Z)
LDFLAGS=-X main.build=$(DATE)
GOGENERATE=$(shell if [ -f .gogenerate ]; then cat .gogenerate; fi)
TESTPACKETS=$(shell if [ -f .testpackages ]; then cat .testpackages; fi)
BENCHPACKETS=$(shell if [ -f .benchpackages ]; then cat .benchpackages; fi)
export GO15VENDOREXPERIMENT=1

PRJ01=$(APP)
BIN01=$(DIR)/bin/$(PRJ01)
VER01=$(shell ${BIN01} version 2>/dev/null)
VERN01=$(shell echo "$(VER01)" | awk -F '-' '{ print $$1 }' )
VERB01=$(shell echo "$(VER01)" | awk -F 'build.' '{ print $$2 }' )

default: dep build

dep:
	mkdir -p ${DIR}/{bin,pkg,run,src,log,conf/keys,doc,tmp,migrations/mysql,storage,template,www} 2>/dev/null; true
	if command -v "gvt"; then cd ${DIR}/src; GOPATH="$(DIR)" gvt update -all; fi
.PHONY: dep

gen:
	for PKGNAME in $(GOGENERATE); do GOPATH="$(DIR)" go generate $${PKGNAME}; done
.PHONY: gen

build:
	GOPATH="$(DIR)" go build -o ${BIN01} -ldflags "${LDFLAGS}" ${PRJ01}
.PHONY: build

dev:
	clear
	${BIN01} --debug daemon
.PHONY: dev

run:
	${BIN01} daemon
.PHONY: run

kill:
	kill -KILL `cat $(DIR)/run/$(PRJ01).pid`
.PHONY: kill

version: v
v:
	${BIN01} version
.PHONY: version
.PHONY: v

gocd:
.PHONY: gocd

## Mysql
m-up:
	@gsmigrate --dir=$(DIR)/migrations/mysql --drv=mysql up
.PHONY: m-up

m-create:
	@gsmigrate --dir=$(DIR)/migrations/mysql --drv=mysql create new
.PHONY: m-create

m-down:
	@gsmigrate --dir=$(DIR)/migrations/mysql --drv=mysql down
.PHONY: m-down

m-status:
	@gsmigrate --dir=$(DIR)/migrations/mysql --drv=mysql status
.PHONY: m-status

test:
	echo "mode: set" > $(DIR)/log/coverage.log
	for PACKET in $(TESTPACKETS); do \
		touch coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=$(DIR)/log/coverage-tmp.log $$PACKET; \
		if [ "$$?" -ne "0" ]; then exit $$?; fi; \
		tail -n +2 $(DIR)/log/coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> $(DIR)/log/coverage.log; \
		rm -f $(DIR)/log/coverage-tmp.log; true; \
	done
.PHONY: test

cover:
	GOPATH=${GOPATH} go tool cover -html=$(DIR)/log/coverage.log
.PHONY: cover

bench:
	for PACKET in $(BENCHPACKETS); do GOPATH=${GOPATH} go test -race -bench=. -benchmem $$PACKET; done
.PHONY: bench

lint:
	gometalinter \
	--vendor \
	--deadline=15m \
	--cyclo-over=20 \
	--disable=aligncheck \
	--skip=src/vendor \
	--skip=src/application/configuration/semver \
	--skip=src/application/configuration/osext \
	--skip=src/application/modules/mimemagic \
	--skip=src/application/models/settings \
	--linter="vet:go tool vet -printf {path}/*.go:PATH:LINE:MESSAGE" \
	src/...
.PHONY: lint

clean:
	rm -rf ${DIR}/bin/*; true
	rm -rf ${DIR}/pkg/*; true
	rm -rf ${DIR}/run/*.pid; true
	rm -rf ${DIR}/log/*.log; true
	rm -rf ${DIR}/rpmbuild; true
	rm -rf ${DIR}/*.log; true
.PHONY: clean
