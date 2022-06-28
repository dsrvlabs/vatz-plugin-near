.PHONY: do_script

build:
	cd bin; /bin/bash build_plugins.sh

start:
	cd bin; /bin/bash start_plugins.sh

stop:
	cd bin; /bin/bash stop_plugins.sh