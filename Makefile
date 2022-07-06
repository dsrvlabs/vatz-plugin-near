ifndef VERBOSE
.SILENT:
endif
SHELL := /bin/bash

build:
	cd bin; $(SHELL) ./build_plugins.sh

start:
	cd bin; $(SHELL) ./start_plugins.sh

stop:
	cd bin; $(SHELL) ./stop_plugins.sh

clean:
	cd bin; $(SHELL) ./clean_plugins.sh