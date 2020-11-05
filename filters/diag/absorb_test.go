package diag

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/ardaguclu/skipper/eskip"
	"github.com/ardaguclu/skipper/filters"
	"github.com/ardaguclu/skipper/logging/loggingtest"
	"github.com/ardaguclu/skipper/proxy/proxytest"
)

const (
	bodySize   = 1 << 12
	logTimeout = 3 * time.Second
)

func testAbsorb(t *testing.T, silent bool) {
	l := loggingtest.New()
	defer l.Close()
	a := withLogger(silent, l)
	fr := make(filters.Registry)
	fr.Register(a)
	p := proxytest.New(
		fr,
		&eskip.Route{
			Filters:     []*eskip.Filter{{Name: "absorb"}},
			BackendType: eskip.ShuntBackend,
		},
	)
	defer p.Close()

	req, err := http.NewRequest(
		"POST",
		p.URL,
		io.LimitReader(
			rand.New(rand.NewSource(time.Now().UnixNano())),
			bodySize,
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Flow-Id", "foo-bar-baz")
	rsp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		t.Fatalf("invalid status code received: %d", rsp.StatusCode)
	}

	expectLog := func(content string, err error) {
		if err != nil {
			t.Fatalf("%s: %v", content, err)
		}
	}

	expectNoLog := func(content string, err error) {
		if err != loggingtest.ErrWaitTimeout {
			t.Fatalf("%s: unexpected log entry", content)
		}
	}

	for _, content := range []string{
		"received request",
		"foo-bar-baz",
		"consumed",
		fmt.Sprint(bodySize),
		"request finished",
	} {
		err := l.WaitFor(content, logTimeout)
		if silent {
			expectNoLog(content, err)
			continue
		}

		expectLog(content, err)
	}
}

func TestAbsorb(t *testing.T) {
	testAbsorb(t, false)
}

func TestAbsorbSilent(t *testing.T) {
	testAbsorb(t, true)
}
