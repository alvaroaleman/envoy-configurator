.PHONY: run
run:
	 ./envoy --v2-config-only -l info -c sample.yaml \
		--service-cluster testcluster \
		--service-node $$(hostname)

main: $(shell find . -name '*.go')
	go build cmd/main.go

run-configurator: main
	./main -cafile CA.crt -crtfile server.crt -keyfile server.key
