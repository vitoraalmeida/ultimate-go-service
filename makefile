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
TELEPRESENCE    := datawire/tel2:2.14.4

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

dev-bill:
	kind load docker-image $(TELEPRESENCE) --name $(KIND_CLUSTER)
	telepresence --context=kind-$(KIND_CLUSTER) helm install
	telepresence --context=kind-$(KIND_CLUSTER) connect

dev-up-local:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml # definição do cluster
	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner
	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)


dev-up: dev-up-local

dev-down-local:
	kind delete cluster --name $(KIND_CLUSTER)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)
	telepresence quit -s
	kind delete cluster --name $(KIND_CLUSTER)

# adiciona imagem no cluster kind para que ele não precise buscar na net
dev-load:
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

# o kustomize gera o texto de descrição dos recursos k8s
# substituindo o que for necessário e passa como argumento
# para o kubectl apply
dev-apply:
	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database
	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --for=condition=Ready
# ------------------------------------------------------------------------------

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide --all-namespaces
	kubectl get pods -o wide --watch --all-namespaces
# ------------------------------------------------------------------------------

# reinicia o deployment em desenvolvimento
dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply

dev-logs:
	#redireciona os logs estruturados que a aplicação gera para a ferramenta de logs legíveis
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(APP)


# ==============================================================================

run-local:
	#redireciona os logs estruturados que a aplicação gera para a ferramenta de logs legíveis
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

run-local-help:
	go run app/services/sales-api/main.go --help

tidy:
	go mod tidy
	go mod vendor

# inicia o serviço de debug
metrics-view-local:
	# expvarmon auxilia na visualização de informações de uso da aplicação
	expvarmon -ports="localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"


test-endpoint:
	curl -il localhost:3000/test

test-endpoint-auth:
	curl -il -H "Authorization: Bearer ${TOKEN}" $(SERVICE_NAME).$(NAMESPACE).svc.cluster.local:3000/test/auth

test-endpoint-auth-local:
	curl -il -H "Authorization: Bearer ${TOKEN}" localhost:3000/test/auth


liveness-local:
	curl -il http://localhost:4000/debug/liveness

liveness:
	curl -il http://$(SERVICE_NAME).$(NAMESPACE).svc.cluster.local:4000/debug/liveness

readiness-local:
	curl -il http://localhost:4000/debug/readiness

readiness:
	curl -il http://$(SERVICE_NAME).$(NAMESPACE).svc.cluster.local:4000/debug/readiness

pgcli-local:
	pgcli postgresql://postgres:postgres@localhost
# token de exemplo gerado com a mesma chave PEM em zark/keys
# ROLE USER: eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiIxMjM0NTY3OCIsImV4cCI6MTcyNjE0ODkzNywiaWF0IjoxNjk0NjEyOTM3LCJSb2xlcyI6WyJVU0VSIl19.qGaka14oavnAbjtcAlGD3Q8orVtRdHFlfUa7LBkC9d3BIKoYkcqav-jnchO-IsaJ27wJXJbS5uwjuBqfpM7bkTJxvYGUIDx7jI3Xp1zmGe0n3pbYJt5nSnucXEi-tCaoU4BhzLAP_f3WV4EeIN9xHY3RTHutbOTm-IdIpln627qIEDrHO1USRBVIZpBLHpaNC7DrlKeCCbPIdh8FfrWFHGUV1SbloxLBfh3RYWzy8swI1LTkVOL6BOLX0Z32nWcSJdpYNH1DmzM65ecI_NaiusoLX9F4QwKyd1WP74ULUoOhghr2zUyJVGPhIXdpd5QEWoBj4zSgCRR0yXqzHSwz_g
# ROLE ADMIN: eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2aWNlIHByb2plY3QiLCJzdWIiOiIxMjM0NTY3OCIsImV4cCI6MTcyNjE0OTA5NSwiaWF0IjoxNjk0NjEzMDk1LCJSb2xlcyI6WyJBRE1JTiJdfQ.G83ASVA98mTSBgPm4QDzEI5Ii0o36fOHNfjlPl01zaqIcteqhRctBa5APxgQWvRmNRB8lG57yfT8_uDFOpH4g6NHVpu1k1OrXbw7chGOIEr6yxlKKzo8T5zOneiNlmBH252f2BCAHio-Fpp-hRl7S3AfXBkAtrwXq0mlIX3gII_BfRCpIqYPrXCpNun9wC3478iReQEMNDf65n2NmexO6q7g70iTOPT1dNZd_iHoaPq7IbsC5OVEY2JX1TcKl59LGwfV1fQsiHb4xgS5xoJsPCpAU7A-6fi4tcYs3WltIzXqH_kogZxpM9yQr3oxvftXaJ2NlqAabECMkSwqCwVHgQ 
