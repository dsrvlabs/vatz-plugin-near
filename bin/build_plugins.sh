#!/bin/bash

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
echo "Build All Plugins"
echo "==================="
for name in "${plugins[@]}"
do
  cd $name
  echo "=> building $name"
  make build >/dev/null
  cd ..
done
echo "==================="
echo "All Build Finished"
