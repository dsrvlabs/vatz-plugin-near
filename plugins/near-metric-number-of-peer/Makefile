SHELL := /bin/bash
build:
	@go build

run:
	$(SHELL) "exec -a near-metric-number-of-peer ./near-metric-number-of-peer >> ../../logs/near-metric-number-of-peer.log 2>&1 &"

exec:
	exec -a "near-metric-number-of-peer" ./near-metric-number-of-peer >> ../../logs/near-metric-number-of-peer.log 2>&1 &

clean:
	rm -rf near-metric-number-of-peer