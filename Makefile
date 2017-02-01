BINARY			= travis-artifacts
NS				= shawnzhu
REPO			= artifacts-v2
BIN_DIR			= bin

SHELL			:= /bin/bash
GOOS			:= $(shell go env GOOS)
GOARCH			:= $(shell go env GOARCH)

KUBECTL     	:= $(shell command -v kubectl 2> /dev/null)
K8S_NAMESPACE 	:= artifacts
K8S_CONTEXT		:= artifacts

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

config_kubectl:
ifndef KUBECTL
	$(error "kubctl is not found please install kubectl")
endif
	echo $(K8S_CERT_B64) | base64 -d > $(HOME)/kube.crt \
	kubectl config set-cluster $(K8S_CLUSTER) --server=$(K8S_SERVER) --certificate-authority=$(HOME)/kube.crt --embed-certs=true \
	kubectl config set-context $(K8S_CONTEXT) --cluster=$(K8S_CLUSTER) --namespace=$(K8S_NAMESPACE) --user=artifacts
	kubectl config set-credentials $(K8S_CLUSTER) --username=$(K8S_USERNAME) --password=$(K8S_PASSWORD)
	kubectl config use-context $(K8S_CONTEXT)

deploy: config_kubectl
	kubectl apply -f k8s-app.yml

clean:
	if [ -d $(BIN_DIR) ]; then rm -r $(BIN_DIR); fi

.PHONY: default test build clean deploy
