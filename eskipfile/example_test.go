package eskipfile_test

import (
	"github.com/ardaguclu/skipper/eskipfile"
	"github.com/ardaguclu/skipper/proxy"
	"github.com/ardaguclu/skipper/routing"
)

func Example() {
	// open file with a routing table:
	dataClient := eskipfile.Watch("/some/path/to/routing-table.eskip")
	defer dataClient.Close()

	// create a routing object:
	rt := routing.New(routing.Options{
		DataClients: []routing.DataClient{dataClient},
	})
	defer rt.Close()

	// create an http.Handler:
	p := proxy.New(rt, proxy.OptionsNone)
	defer p.Close()
}
