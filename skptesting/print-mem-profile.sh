#! /bin/bash

go tool pprof -text $GOPATH/src/github.com/ardaguclu/skipper/skptesting/mem-profile.prof
