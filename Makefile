BINARY			= travis-artifacts
NS				= shawnzhu
REPO			= artifacts-v2
BIN_DIR			= bin

SHELL			:= /bin/bash
GOOS			:= $(shell go env GOOS)
GOARCH			:= $(shell go env GOARCH)

KUBECTL     	:= $(shell command -v kubectl 2> /dev/null)
KUBECTL_URL		:= https://storage.googleapis.com/kubernetes-release/release/v1.5.2/bin/linux/amd64/kubectl
K8S_NAMESPACE 	:= artifacts
K8S_CONTEXT		:= artifacts
K8S_USERNAME	:= travisci

default: clean install test build

install:
	go get -t -v ./...

test:
ifeq ($(TRAVIS),true)
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	goveralls -v
else
	go test -v ./...
endif

build:
	test -d $(BIN_DIR) || mkdir -p $(BIN_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BINARY) cmd/travis-artifacts/main.go

release:
	docker build -t $(NS)/$(REPO):$(TAG) .
	docker push $(NS)/$(REPO):$(TAG)

install_kubectl:
ifdef KUBECTL
	$(info "kubectl is installed already")
else
	wget -O /usr/local/bin/kubectl -nv $(KUBECTL_URL)
	chmod +x /usr/local/bin/kubectl
endif

config_kubectl: install_kubectl
ifdef KUBECONFIG
	$(info "kubectl is configured already")
else
	@printf $(K8S_CERT_B64) | base64 -d > $(HOME)/kube.crt
	@kubectl config set-credentials $(K8S_USERNAME) --token=$(shell printf $$K8S_TOKEN_B64 | base64 -d)
	kubectl config set-cluster $(K8S_CLUSTER) --server=$(K8S_SERVER) --certificate-authority=$(HOME)/kube.crt --embed-certs=true
	kubectl config set-context $(K8S_CONTEXT) --cluster=$(K8S_CLUSTER) --namespace=$(K8S_NAMESPACE) --user=$(K8S_USERNAME)
	kubectl config use-context $(K8S_CONTEXT)
endif

deploy: config_kubectl
	kubectl apply -f k8s-app.yml

clean:
	if [ -d $(BIN_DIR) ]; then rm -r $(BIN_DIR); fi

.PHONY: default test build clean deploy
