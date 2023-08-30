# Verifica se podemos usar o ash em Alpine images ou altera o default para BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

run-local:
	go run app/services/sales-api/main.go

tidy:
	go mod tidy
	go mod vendor
