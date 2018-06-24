.PHONY: run
run:
	 ./envoy --v2-config-only -l info -c sample.yaml \
		--service-cluster testcluster \
		--service-node $$(hostname)
