.PHONY: run
run:
	 ./envoy --v2-config-only -l info -c sample.yaml \
		--service-cluster testcluster \
		--service-node $$(hostname)

main: cmd/main.go pkg/controller/*.go
	go build cmd/main.go

run-configurator: main
	./main
