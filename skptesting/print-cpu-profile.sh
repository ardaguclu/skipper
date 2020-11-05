#! /bin/bash

go tool pprof -text $GOPATH/bin/skptesting $GOPATH/src/github.com/ardaguclu/skipper/skptesting/cpu-profile.prof
