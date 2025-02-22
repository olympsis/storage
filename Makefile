VERSION := v0.2.0
SERVICE_NAME := storage

REPO_NAME := main
PROJECT_ID := olympsis-408521
LOCATION := us-central1-docker.pkg.dev

.PHONY: all dep build clean test coverage coverhtml lint

all: build

build: dep ## Build the binary file
	go build

local:
	docker build . -t $(SERVICE_NAME)
	docker run -p 80:80 $(SERVICE_NAME):latest

publish:
	docker build . -t $(SERVICE_NAME) --platform linux/amd64 --build-arg VERSION=$(VERSION)
	docker tag $(SERVICE_NAME) $(LOCATION)/$(PROJECT_ID)/$(REPO_NAME)/$(SERVICE_NAME):$(VERSION)
	docker push $(LOCATION)/$(PROJECT_ID)/$(REPO_NAME)/$(SERVICE_NAME):$(VERSION)

update-service: #Updates the linux service
	make build && \
	if [ $$? -ne 0 ]; then \
		echo "Error: Failed to build new server binary." && \
		exit 1; \
	fi && \
	rm /sbin/olympsis-storage && \
	mv storage /sbin/olympsis-storage && \
	if [ $$? -ne 0 ]; then \
		echo "Error: Failed to move binary." && \
		exit 1; \
	fi && \
	systemctl restart olympsis-storage.service && \
	echo "Update Successful"