SHELL := /bin/bash
build:
	@go build

run:
	$(SHELL) -c "exec -a machine-status-cpu ./machine-status-cpu >> ../../logs/machine-status-cpu.log 2>&1 &"

exec:
	exec -a "machine-status-cpu" ./machine-status-cpu >> ../../logs/machine-status-cpu.log 2>&1 &

clean:
	rm -rf machine-status-cpu