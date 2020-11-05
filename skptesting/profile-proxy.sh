#! /bin/bash

if [ "$1" == -help ]; then
	log profile-proxy.sh [duration] [connections]
	exit 0
fi

source $GOPATH/src/github.com/ardaguclu/skipper/skptesting/benchmark.inc

trap cleanup SIGINT

log [generating content]
lorem
log [content generated]

log; log [starting servers]
# ngx nginx-static.conf
skp :9980 static.eskip
skp-pprof :9090 proxy.eskip
log [servers started, wait 1 sec]
sleep 1

log; log [profiling skipper]
bench :9090
log [profiling skipper done]

cleanup
log; log [all done]
