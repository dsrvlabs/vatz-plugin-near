#!/bin/bash

## declare an array variable
declare -a plugins=("machine-status-cpu"
                "machine-status-disk"
                "machine-status-memory"
                "near-metric-alive"
                "near-metric-block-height"
                "near-metric-chunk-produce-rate"
                "near-metric-number-of-peer"
                "near-metric-uptime"
                )

cd ..
cd plugins
echo "Cleaning All Plugins"
echo "==================="
for name in "${plugins[@]}"
do
  cd $name
  echo "=> cleaning $name"
  make clean >/dev/null
  cd ..
done
echo "==================="
echo "All Plugins Cleaned"

