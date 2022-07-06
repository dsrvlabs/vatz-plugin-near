SHELL := /bin/bash
build:
	@go build

run:
	$(SHELL) "exec -a machine-status-memory ./machine-status-memory >> ../../logs/machine-status-memory.log 2>&1 &"

exec:
	exec -a "machine-status-memory" ./machine-status-memory >> ../../logs/machine-status-memory.log 2>&1 &
clean:
	rm -rf machine-status-memory
