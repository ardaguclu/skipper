#! /bin/bash

cd $GOPATH/src/github.com/ardaguclu/skipper/skptesting
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -nodes -days 9
