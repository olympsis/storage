NAME := storage
VERSION := v0.2.0

REPO_NAME := main
PROJECT_ID := olympsis-408521
LOCATION := us-central1-docker.pkg.dev

.PHONY: all dep build clean test coverage coverhtml lint

all: build

build: dep ## Build the binary file
	go build

docker:
	$(MAKE) docker-build
	$(MAKE) docker-run

docker-build:
	docker build . -f ./tools/Dockerfile -t $(NAME)

docker-run:
	docker run -p 7002:80 $(NAME):latest

docker-publish:
	docker build . -t $(NAME) --platform linux/amd64 --build-arg VERSION=$(VERSION)
	docker tag $(NAME) $(LOCATION)/$(PROJECT_ID)/$(REPO_NAME)/$(NAME):$(VERSION)
	docker push $(LOCATION)/$(PROJECT_ID)/$(REPO_NAME)/$(NAME):$(VERSION)

update-service: ## Updates the linux service
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