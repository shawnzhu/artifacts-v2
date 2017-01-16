BINARY		= travis-artifacts
NS			= shawnzhu
REPO		= artifacts-v2
BIN_DIR		= bin

SHELL		:= /bin/bash
GOOS		:= $(shell go env GOOS)
GOARCH		:= $(shell go env GOARCH)


default: clean install test build

install:
	go get -t -v ./...

test:
	if [ "${TRAVIS}" == "true" ]; then		\
		go get golang.org/x/tools/cmd/cover;\
		go get github.com/mattn/goveralls;	\
		goveralls -v;						\
	else									\
		go test -v ./...;					\
	fi

build:
	test -d $(BIN_DIR) || mkdir -p $(BIN_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BINARY) cmd/travis-artifacts/main.go

release:
	docker build -t $(NS)/$(REPO):$(TAG) .
	docker push $(NS)/$(REPO):$(TAG)

install_kubectl:
	if [ "${TRAVIS}" == "true" ] && [ "$(which kubectl)" == "" ]; then		\
		KUBECTL_RELEASE="$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)" \
		sudo curl -ssL -o /usr/local/bin/kubectl \
		https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_RELEASE}/bin/linux/amd64/kubectl \
		sudo chmod +x /usr/local/bin/kubectl \
	fi

deploy:
	/usr/local/bin/kubectl apply -f k8s-app.yml

clean:
	if [ -d $(BIN_DIR) ]; then rm -r $(BIN_DIR); fi

.PHONY: default test build clean deploy
