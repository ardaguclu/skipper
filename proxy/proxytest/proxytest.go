package proxytest

import (
	"net/http/httptest"
	"time"

	"github.com/ardaguclu/skipper/eskip"
	"github.com/ardaguclu/skipper/filters"
	"github.com/ardaguclu/skipper/loadbalancer"
	"github.com/ardaguclu/skipper/logging/loggingtest"
	"github.com/ardaguclu/skipper/proxy"
	"github.com/ardaguclu/skipper/routing"
	"github.com/ardaguclu/skipper/routing/testdataclient"
)

type TestProxy struct {
	URL     string
	Log     *loggingtest.Logger
	routing *routing.Routing
	proxy   *proxy.Proxy
	server  *httptest.Server
}

func WithRoutingOptions(fr filters.Registry, o routing.Options, routes ...*eskip.Route) *TestProxy {
	return newTestProxy(fr, o, proxy.Params{CloseIdleConnsPeriod: -time.Second}, routes...)
}

func WithParams(fr filters.Registry, proxyParams proxy.Params, routes ...*eskip.Route) *TestProxy {
	return newTestProxy(fr, routing.Options{}, proxyParams, routes...)
}

func newTestProxy(fr filters.Registry, routingOptions routing.Options, proxyParams proxy.Params, routes ...*eskip.Route) *TestProxy {
	tl := loggingtest.New()

	if len(routingOptions.DataClients) == 0 {
		dc := testdataclient.New(routes)
		routingOptions.DataClients = []routing.DataClient{dc}
	}

	routingOptions.FilterRegistry = fr
	routingOptions.Log = tl
	routingOptions.PostProcessors = []routing.PostProcessor{loadbalancer.NewAlgorithmProvider()}

	rt := routing.New(routingOptions)
	proxyParams.Routing = rt

	pr := proxy.WithParams(proxyParams)
	tsp := httptest.NewServer(pr)

	if err := tl.WaitFor("route settings applied", 3*time.Second); err != nil {
		panic(err)
	}

	return &TestProxy{
		URL:     tsp.URL,
		Log:     tl,
		routing: rt,
		proxy:   pr,
		server:  tsp,
	}
}

func New(fr filters.Registry, routes ...*eskip.Route) *TestProxy {
	return WithParams(fr, proxy.Params{CloseIdleConnsPeriod: -time.Second}, routes...)
}

func (p *TestProxy) Close() error {
	p.Log.Close()
	p.routing.Close()

	err := p.proxy.Close()
	if err != nil {
		return err
	}

	p.server.Close()
	return nil
}
