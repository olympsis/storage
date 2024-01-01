VERSION 	 := v0.2.0
SERVICE_NAME := olympsis/storage
PKG := "$(SERVICE_NAME)"
PKG_LIST := $( go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $( find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

export PORT=80

.PHONY: all dep build clean test coverage coverhtml lint

all: build

lint: ## Lint the files
	golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	go test -short ${PKG_LIST}

race: dep ## Run data race detector
	go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	go test -msan -short ${PKG_LIST}

dep: ## Get the dependencies
	go get -v -d ./...

build: dep ## Build the binary file
	go build -v $(PKG) 

docker:
	docker build . -t $(SERVICE_NAME) --platform linux/amd64 --build-arg VERSION=$(VERSION)
	docker tag $(SERVICE_NAME):latest $(SERVICE_NAME):$(VERSION)
	docker push $(SERVICE_NAME):$(VERSION)

local:
	docker build . -t $(SERVICE_NAME)
	docker run -p 8082:8082 $(SERVICE_NAME):latest

run:
	go run -x main.go

clean: ## Remove previous build
	rm -f $(SERVICE_NAME)