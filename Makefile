VERSION := v0.2.0
SERVICE_NAME := storage

REPO_NAME := main
PROJECT_ID := olympsis-408521
LOCATION := us-central1-docker.pkg.dev

.PHONY: all dep build clean test coverage coverhtml lint

all: build

local:
	docker build . -t $(SERVICE_NAME)
	docker run -p 80:80 $(SERVICE_NAME):latest

publish:
	docker build . -t $(SERVICE_NAME) --platform linux/amd64 --build-arg VERSION=$(VERSION)
	docker tag $(SERVICE_NAME) $(LOCATION)/$(PROJECT_ID)/$(REPO_NAME)/$(SERVICE_NAME):$(VERSION)
	docker push $(LOCATION)/$(PROJECT_ID)/$(REPO_NAME)/$(SERVICE_NAME):$(VERSION)