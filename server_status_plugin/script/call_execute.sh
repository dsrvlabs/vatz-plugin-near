#!/bin/bash

grpcurl -plaintext \
    -d "{\"execute_info\": {\"function\": \"$1\"}}" \
    localhost:9091 \
    pilot.plugin.ManagerPlugin.execute
