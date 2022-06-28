#!/bin/bash

## declare an array variable
declare -a plugins=("machine-status-cpu"
                "machine-status-disk"
                "machine-status-memory"
                "near-metric-blockheight"
                "near-metric-up"
                )


for i in "${plugins[@]}"
do
   PID=`ps -eaf | grep $i | grep -v grep | awk '{print $2}'`
   if [[ "" !=  "$PID" ]]; then
     echo "Stopping Plugins: $i in PID: $PID"
     kill -15 $PID
   fi
done

