package main

import (
	"github.com/ardaguclu/skipper/eskip"
	"github.com/ardaguclu/skipper/routing"
)

type DataClient string

func InitDataClient([]string) (routing.DataClient, error) {
	var dc DataClient = ""
	return dc, nil
}

func (dc DataClient) LoadAll() ([]*eskip.Route, error) {
	return eskip.Parse(string(dc))
}

func (dc DataClient) LoadUpdate() ([]*eskip.Route, []string, error) {
	return nil, nil, nil
}
