#!/bin/bash

kill -15 $(lsof -t -i:9091)
kill -15 $(lsof -t -i:9092)
