SHELL := /bin/bash
build:
	@go build

run:
	$(SHELL) -c "exec -a near-metric-block-height ./near-metric-block-height >> ../../logs/near-metric-block-height.log 2>&1 &"

exec:
	exec -a "near-metric-block-height" ./near-metric-block-height >> ../../logs/near-metric-block-height.log 2>&1 &

clean:
	rm -rf near-metric-block-height
