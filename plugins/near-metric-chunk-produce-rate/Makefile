SHELL := /bin/bash
build:
	@go build

run:
	$(SHELL) "exec -a near-metric-chunk-produce-rate ./near-metric-chunk-produce-rate >> ../../logs/near-metric-chunk-produce-rate.log 2>&1 &"

exec:
	exec -a "near-metric-chunk-produce-rate" ./near-metric-chunk-produce-rate >> ../../logs/near-metric-chunk-produce-rate.log 2>&1 &

clean:
	rm -rf near-metric-chunk-produce-rate