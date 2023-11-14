#!/usr/bin/env bash

mkdir -p `pwd`/logs
ofile=`pwd`/logs/access_`date +%F_%H%M%S`.log
echo Logging into: $ofile
unbuffer time ./example.exe rest | tee $ofile
