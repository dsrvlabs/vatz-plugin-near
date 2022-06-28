#!/bin/bash

declare -a plugins=("machine-status-cpu"
                "machine-status-disk"
                "machine-status-memory"
                "near-metric-blockheight"
                "near-metric-up"
                )

cd ..
cd plugins

for name in "${plugins[@]}"
do
  cd $name
  echo "build $name"
  make build
  cd ..
done
