SHELL := /bin/bash
build:
	@go build

run:
	$(SHELL) "exec -a near-metric-alive ./near-metric-alive >> ../../logs/near-metric-alive.log 2>&1 &"

exec:
	exec -a "near-metric-alive" ./near-metric-alive >> ../../logs/near-metric-alive.log 2>&1 &

clean:
	rm -rf near-metric-alive
