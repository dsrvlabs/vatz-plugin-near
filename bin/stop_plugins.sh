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

echo "Stopping All Plugins"
echo "==================="
for i in "${plugins[@]}"
do
   PID=`ps -eaf | grep $i | grep -v grep | awk '{print $2}'`
   if [[ "" !=  "$PID" ]]; then
     echo "=> Stopping Plugins: $i in PID: $PID"
     kill -15 $PID >/dev/null
   fi
done
echo "==================="
echo "All Plugins has stopped"

