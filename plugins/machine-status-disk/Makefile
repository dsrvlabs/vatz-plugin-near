SHELL := /bin/bash
build:
	@go build

run:
	$(SHELL) "exec -a machine-status-disk ./machine-status-disk >> ../../logs/machine-status-disk.log 2>&1 &"

exec:
	exec -a "machine-status-disk" ./machine-status-disk >> ../../logs/machine-status-disk.log 2>&1 &

clean:
	rm -rf machine-status-disk