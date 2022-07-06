SHELL := /bin/bash
build:
	@go build

run:
	$(SHELL) -c "exec -a near-metric-uptime ./near-metric-uptime --network testnet >> ../../logs/near-metric-uptime.log 2>&1 &"

exec:
	exec -a "near-metric-uptime" ./near-metric-uptime >> ../../logs/near-metric-uptime.log 2>&1 &

clean:
	rm -rf near-metric-uptime