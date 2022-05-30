APP = metida
SERVICE_PATH = github.com/Dsmit05
BR = `git rev-parse --symbolic-full-name --abbrev-ref HEAD`
VER = `git describe --tags --abbrev=0`
TIMESTM = `date -u '+%Y-%m-%d_%H:%M:%S%p'`
FORMAT = $(VER)-$(TIMESTM)
DOCTAG = $(VER)-$(BR)

.PHONY: info
info:
	make -v
	go version
	echo "appname:"$(APP) "version:"$(FORMAT) "docker_tag:"$(DOCTAG)

.PHONY: unit-test
unit-test:
	CGO_ENABLED=0 go test ./...


.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(APP) -ldflags "-X $(SERVICE_PATH)/$(APP)/internal/config.buildVersion=$(FORMAT)"


.PHONY: build-image
build-image:
	sudo docker build -t $(APP):$(DOCTAG) .


.PHONY: run-app
run-app:
	docker run -d --name=$(APP)-$(VER) -p 8080:8080 $(APP):$(DOCTAG)

.PHONY: del-app
del-app:
	docker rm $(APP)-$(VER)

.PHONY: swag-init
swag-init:
	swag init

.PHONY: start
start:
	docker-compose up -d

.PHONY: recreate-api
recreate-api:
	docker-compose up -d --force-recreate --no-deps --build api

.PHONY: sql-generate
sql-generate:
	sqlc generate


