## Simple projects tooling for every day
## (c)Alex Geer <monoflash@gmail.com>
## Makefile version: 16.09.2019

## Project name and source directory path
export DIR  := $(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

## Creating .env file from template, if file not exists
ifeq ("$(wildcard $(DIR)/.env)","")
  RSP1      := $(shell cp -v $(DIR)/.example_env $(DIR)/.env)
endif
## Creating .prj file from template, if file not exists
ifeq ("$(wildcard $(DIR)/.prj)","")
  RSP2      := $(shell cp -v $(DIR)/.example_prj $(DIR)/.prj)
endif
include $(DIR)/.env
include $(DIR)/.prj

APP         := $(PROJECT_NAME)
GOPATH      := $(DIR):$(GOPATH)
DATE        := $(shell date -u +%Y%m%d.%H%M%S.%Z)
LDFLAGS      = -X main.build=$(DATE) $(PRJ_LDFLAGS:'')
GOGENERATE   = $(shell if [ -f .gogenerate ]; then cat .gogenerate; fi)
TESTPACKETS  = $(shell if [ -f .testpackages ]; then cat .testpackages; fi)
BENCHPACKETS = $(shell if [ -f .benchpackages ]; then cat .benchpackages; fi)
GO111MODULE ?= $(GO111MODULE:off)

PRJ01       := $(APP)
BIN01       := $(DIR)/bin/$(PRJ01)
VER01       := $(shell ${BIN01} version 2>/dev/null)
VERN01      := $(shell echo "$(VER01)" | awk -F '-' '{ print $$1 }' )
VERB01      := $(shell echo "$(VER01)" | awk -F 'build.' '{ print $$2 }' )
PIDF01      := $(DIR)/run/$(PRJ01).pid
PIDN01       = $(shell if [ -f $(PIDF01) ]; then  cat $(PIDF01); fi)

PRJ_CGO_ENABLED = $(PRJ_CGO_ENABLED:0)

default: help

## Dependences manager
dep-init:
	@for dir in ${PROJECT_FOLDERS}; do \
	  if [ ! -d "${DIR}/$${dir}" ]; then \
		  mkdir -p "${DIR}/$${dir}"; \
		fi; \
	done
.PHONY: dep-init
dep: dep-init
	@for cmd in ${PROJECT_DEPENDENCES}; do bash -c "$${cmd}"; done
.PHONY: dep

## Code generation (run only during development)
# All generating files are included in a .gogenerate file
gen: dep-init
	@for PKGNAME in $(GOGENERATE); do GOPATH="$(DIR)" DB2STRUCT_DRV="$(GOOSE_DRV_MYSQL)" DB2STRUCT_DSN="$(GOOSE_DSN_MYSQL)" go generate $${PKGNAME}; done
.PHONY: gen

## Build project
build:
	@GO111MODULE="off" GOPATH="$(DIR)" CGO_ENABLED=$(PRJ_CGO_ENABLED) go build -a -i \
	-o ${BIN01} \
	-gcflags "all=-N -l" \
	-ldflags "${LDFLAGS}" \
	-pkgdir ${DIR}/pkg \
	${PRJ01}
.PHONY: build

## Build project for i386 architecture
build-i386:
	@GO111MODULE="off" GOPATH="$(DIR)" CGO_ENABLED=$(PRJ_CGO_ENABLED) GOARCH=386 go build -a -i \
	-o ${BIN01} \
	-gcflags "all=-N -l" \
	-ldflags "${LDFLAGS}" \
	-pkgdir ${DIR}/pkg \
	${PRJ01}
.PHONY: build-i386

## Run application in development mode
dev: clear
	@for cmd in $(PROJECT_RUN_DEVELOPMENT); do cd ${DIR}; sh -c "$${cmd}"; done
.PHONY: dev

## Run application in production mode
run:
	@for cmd in $(PROJECT_RUN_PRODUCTION); do cd ${DIR}; sh -c "$${cmd}"; done
.PHONY: run

## Kill process and remove pid file
kill:
	@if [ ! "$(PIDN01)x" == "x" ]; then \
		kill -KILL "$(PIDN01)" 2>/dev/null; \
		if [ $$? -ne 0 ]; then echo "No such process ID: $(PIDN01)"; fi; \
	fi
	@rm "$(PIDF01)" 2>/dev/null; true
.PHONY: kill

## Getting application version
version: v
v:
	@${BIN01} version
.PHONY: version
.PHONY: v

## RPM build openSUSE linux version
RPMBUILD_OS ?= $(RPMBUILD_OS:leap)
RPMBUILD_OS ?= $(RPMBUILD_OS:tumbleweed)
## Creating RPM package
rpm:
	## Prepare for creating RPM package
	@mkdir -p ${DIR}/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}; true
	## Copying the content needed to build the RPM package
	# File descriptions are contained in the .rpm file
	@for item in $(RPM_BUILD_SOURCE); do\
		SRC=`echo $${item} | awk -F':' '{print $$1}'`; \
		DST=`echo $${item} | awk -F':' '{print $$2}'`; \
		cp -v ${DIR}/$${SRC} ${DIR}/rpmbuild/$${DST}; \
	done
	## Execution of data preparation commands for build an RPM package
	# Command descriptions are contained in the .rpm file
	@for cmd in $(RPM_BUILD_COMMANDS); do\
		cd ${DIR}; sh -v -c "$${cmd}"; \
	done
	## Updates SPEC changelog section, from git log information
	@if command -v "changelogmaker"; then \
		mv ${DIR}/rpmbuild/SPECS/${PRJ01}.spec ${DIR}/rpmbuild/SPECS/src.spec; \
		cd ${DIR}; changelogmaker -s ${DIR}/rpmbuild/SPECS/src.spec > ${DIR}/rpmbuild/SPECS/${PRJ01}.spec; \
	fi
	## Build the RPM package
	@RPMBUILD_OS="${RPMBUILD_OS}" rpmbuild \
	  --target x86_64 \
		--define "_topdir ${DIR}/rpmbuild" \
	  	--define "_app_version_number $(VERN01)" \
	  	--define "_app_version_build $(VERB01)" \
	  	-bb ${DIR}/rpmbuild/SPECS/${PRJ01}.spec
.PHONY: rpm

## Creating RPM package for i386 architecture
rpm-i386:
	## Prepare for creating RPM package
	@mkdir -p ${DIR}/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}; true
	## Copying the content needed to build the RPM package
	# File descriptions are contained in the .rpm file
	@for item in $(RPM_BUILD_SOURCE_I386); do\
		SRC=`echo $${item} | awk -F':' '{print $$1}'`; \
		DST=`echo $${item} | awk -F':' '{print $$2}'`; \
		cp -v ${DIR}/$${SRC} ${DIR}/rpmbuild/$${DST}; \
	done
	## Execution of data preparation commands for build an RPM package
	# Command descriptions are contained in the .rpm file
	@for cmd in $(RPM_BUILD_COMMANDS_I386); do\
		cd ${DIR}; sh -v -c "$${cmd}"; \
	done
	## Updates SPEC changelog section, from git log information
	@if command -v "changelogmaker"; then \
		mv ${DIR}/rpmbuild/SPECS/${PRJ01}.spec ${DIR}/rpmbuild/SPECS/src.spec; \
		cd ${DIR}; changelogmaker -s ${DIR}/rpmbuild/SPECS/src.spec > ${DIR}/rpmbuild/SPECS/${PRJ01}.spec; \
	fi
	## Build the RPM package
	@RPMBUILD_OS="${RPMBUILD_OS}" rpmbuild \
	  --target i386 \
		--define "_topdir ${DIR}/rpmbuild" \
	  	--define "_app_version_number $(VERN01)" \
	  	--define "_app_version_build $(VERB01)" \
	  	-bb ${DIR}/rpmbuild/SPECS/${PRJ01}.spec
.PHONY: rpm-i386

## Migration tools for all databases
# Please see files .env and .env_example, for setup access to databases
####################################
COMMANDS  = up create down status redo version
MTARGETS := $(shell \
for cmd in $(COMMANDS); do \
	for drv in $(MIGRATIONS); do \
		echo "m-$${drv}-$${cmd}"; \
	done; \
done)
## Migration tools create directory
migration-mkdir:
	@for dir in $$(echo $(MIGRATIONS)); do \
		mkdir -p "$(DIR)/migrations/$${dir}"; true; \
	done
.PHONY: migration-mkdir
## Migration tools gets data from env
MIGRATION_DIR  = ${$(shell echo $(shell echo "${@}" | sed -e 's/^m-\(.*\)-\(.*\)$$/\1/') | awk '{print "GOOSE_DIR_"toupper($$0)}')}
MIGRATION_DRV  = ${$(shell echo $(shell echo "${@}" | sed -e 's/^m-\(.*\)-\(.*\)$$/\1/') | awk '{print "GOOSE_DRV_"toupper($$0)}')}
MIGRATION_DSN  = ${$(shell echo $(shell echo "${@}" | sed -e 's/^m-\(.*\)-\(.*\)$$/\1/') | awk '{print "GOOSE_DSN_"toupper($$0)}')}
MIGRATION_CMD  = $(shell echo $(shell echo "${@}" | sed -e 's/^m-\(.*\)-\(.*\)$$/\2/'))
MIGRATION_TMP := $(shell mktemp)
## Migration tools targets
migration-commands: $(MTARGETS)
$(MTARGETS): migration-mkdir
	@if [ "$(MIGRATION_CMD)" == "create" ]; then\
		read -p "Введите название миграции: " MGRNAME && \
		if [ "$${MGRNAME}" == "" ]; then MGRNAME="new"; fi && \
		echo "$${MGRNAME}" > "$(MIGRATION_TMP)"; \
	fi
	@if ([ ! "`cat $(MIGRATION_TMP)`" = "" ]) && ([ "$(MIGRATION_CMD)" == "create" ]); then\
		GOOSE_DIR="$(MIGRATION_DIR)" GOOSE_DRV="$(MIGRATION_DRV)" GOOSE_DSN="$(MIGRATION_DSN)" gsmigrate $(MIGRATION_CMD) "`cat $(MIGRATION_TMP)`"; \
	else \
		GOOSE_DIR="$(MIGRATION_DIR)" GOOSE_DRV="$(MIGRATION_DRV)" GOOSE_DSN="$(MIGRATION_DSN)" gsmigrate $(MIGRATION_CMD); \
	fi
	@if [ -f "$(MIGRATION_TMP)" ]; then rm "$(MIGRATION_TMP)"; fi
.PHONY: migration-commands $(MTARGETS)
####################################

## Testing one or multiple packages as well as applications with reporting on the percentage of test coverage
# All testing files are included in a .testpackages file
test:
	@echo "mode: set" > $(DIR)/log/coverage.log
	@for PACKET in $(TESTPACKETS); do \
		touch coverage-tmp.log; \
		GOPATH=${GOPATH} go test -v -covermode=count -coverprofile=$(DIR)/log/coverage-tmp.log $$PACKET; \
		if [ "$$?" -ne "0" ]; then exit $$?; fi; \
		tail -n +2 $(DIR)/log/coverage-tmp.log | sort -r | awk '{if($$1 != last) {print $$0;last=$$1}}' >> $(DIR)/log/coverage.log; \
		rm -f $(DIR)/log/coverage-tmp.log; true; \
	done
.PHONY: test

## Displaying in the browser coverage of tested code, on the html report (run only during development)
cover: test
	@GOPATH=${GOPATH} go tool cover -html=$(DIR)/log/coverage.log
.PHONY: cover

## Performance testing
# All testing files are included in a .benchpackages file
bench:
	@for PACKET in $(BENCHPACKETS); do GOPATH=${GOPATH} go test -race -bench=. -benchmem $$PACKET; done
.PHONY: bench

## Code quality testing
# https://github.com/alecthomas/gometalinter/
# install: curl -L https://git.io/vp6lP | sh
lint:
	@golangci-lint \
	run \
	--enable-all \
	--disable nakedret \
	--disable gochecknoinits \
	src/...
.PHONY: lint

## Cleaning console screen
clear:
	clear
.PHONY: clear

## Clearing project temporary files
clean:
	@GOPATH="$(DIR)" go clean -cache
	@chown -R `whoami` ${DIR}/pkg/; true
	@chmod -R 0777 ${DIR}/pkg/; true
	@rm -rf ${DIR}/bin/*; true
	@rm -rf ${DIR}/pkg/*; true
	@rm -rf ${DIR}/run/*.pid; true
	@rm -rf ${DIR}/log/*.log; true
	@rm -rf ${DIR}/rpmbuild; true
	@rm -rf ${DIR}/*.log; true
	@export DIR=
.PHONY: clean

## Help for main targets
help:
	@echo "Usage: make [target]"
	@echo "  target is:"
	@echo "    dep                  - Загрузка и одновление зависимостей проекта"
	@echo "    gen                  - Кодогенерация с использованием go generate"
	@echo "    build                - Компиляция приложения"
	@echo "    build-i386           - Компиляция приложения для архитектуры i386"
	@echo "    run                  - Запуск приложения в продакшн режиме"
	@echo "    dev                  - Запуск приложения в режиме разработки"
	@echo "    kill                 - Отправка приложению сигнала kill -HUP, используется в случае зависания"
	@echo "    m-[driver]-[command] - Работа с миграциями базы данных"
	@echo "                           Используемые базы данных (driver) описываются в файле .env"
	@echo "                           Доступные драйвера баз данных: mysql clickhouse sqlite3 postgres redshift tidb"
	@echo "                           Доступные команды: up, down, create, status, redo, version"
	@echo "                           Пример команд при включённой базе данных mysql:"
	@echo "                             make m-mysql-up      - примернение миграций до самой последней версии"
	@echo "                             make m-mysql-down    - отмена последней миграции"
	@echo "                             make m-mysql-create  - создание нового файла миграции"
	@echo "                             make m-mysql-status  - статус всех миграций базы данных"
	@echo "                             make m-mysql-redo    - отмена и повторное применение последней миграции"
	@echo "                             make m-mysql-version - отображение версии базы данных (применённой миграции)"
	@echo "                           Подробная информаци по командам доступна в документации утилиты gsmigrate"
	@echo "    version              - Вывод на экран версии приложения"
	@echo "    rpm                  - Создание RPM пакета"
	@echo "    rpm-i386             - Создание RPM пакета для архитектуры i386"
	@echo "    bench                - Запуск тестов производительности проекта"
	@echo "    test                 - Запуск тестов проекта"
	@echo "    cover                - Запуск тестов проекта с отображением процента покрытия кода тестами"
	@echo "    lint                 - Запуск проверки кода с помощью gometalinter"
	@echo "    clean                - Очистка папки проекта от временных файлов"
.PHONY: help
