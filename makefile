# Verifica se podemos usar o ash em Alpine images ou altera o default para BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# CLASS NOTES
#
# Kind
# 	For full Kind v0.18 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.18.0

# ==============================================================================
# Define versões das imagens docker utilizadas

GOLANG          := golang:1.20
ALPINE          := alpine:3.18
KIND            := kindest/node:v1.27.1
POSTGRES        := postgres:15.3
VAULT           := hashicorp/vault:1.13
ZIPKIN          := openzipkin/zipkin:2.24 # telemetria
TELEPRESENCE    := datawire/tel2:2.13.1

KIND_CLUSTER    := go-dev-cluster
NAMESPACE       := sales-system
APP             := sales
BASE_IMAGE_NAME := vitoraalmeida/service
SERVICE_NAME    := sales-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME)-metrics:$(VERSION)

# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"

# ==============================================================================
# Constrói containers

all: service

service:
	docker build \
		-f zarf/docker/dockerfile.service \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Utilizando KIND - Kubernetes On Docker

dev-up-local:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml # definição do cluster

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-up: dev-up-local

dev-down-local:
	kind delete cluster --name $(KIND_CLUSTER)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

# ------------------------------------------------------------------------------

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

# ==============================================================================

run-local:
	go run app/services/sales-api/main.go

tidy:
	go mod tidy
	go mod vendor
