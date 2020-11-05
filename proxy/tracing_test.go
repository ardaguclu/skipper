package proxy

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ardaguclu/skipper/tracing/tracingtest"

	"github.com/opentracing/opentracing-go/mocktracer"
)

const traceHeader = "X-Trace-Header"

func TestTracingFromWire(t *testing.T) {
	traceContent := fmt.Sprintf("%x", md5.New().Sum([]byte(time.Now().String())))
	s := startTestServer(nil, 0, func(r *http.Request) {
		th, ok := r.Header[traceHeader]
		if !ok {
			t.Errorf("missing %s request header", traceHeader)
		} else {
			if th[0] != traceContent {
				t.Errorf("wrong X-Trace-Header content: %s", th[0])
			}
		}
	})
	defer s.Close()

	u, _ := url.ParseRequestURI("https://www.example.org/hello")
	r := &http.Request{
		URL:    u,
		Method: "GET",
		Header: make(http.Header),
	}
	r.Header.Set(traceHeader, traceContent)
	w := httptest.NewRecorder()

	doc := fmt.Sprintf(`hello: Path("/hello") -> "%s"`, s.URL)
	tracer := &tracingtest.Tracer{}
	params := Params{
		OpenTracing: &OpenTracingParams{
			Tracer: tracer,
		},
		Flags: FlagsNone,
	}

	tp, err := newTestProxyWithParams(doc, params)
	if err != nil {
		t.Error(err)
		return
	}
	defer tp.close()

	tp.proxy.ServeHTTP(w, r)

	if len(tracer.RecordedSpans) == 0 {
		t.Fatal("no span recorded...")
	}
	if tracer.RecordedSpans[0].Trace != traceContent {
		t.Errorf("trace not found, got `%s` instead", tracer.RecordedSpans[0].Trace)
	}
	if len(tracer.RecordedSpans[0].Refs) == 0 {
		t.Errorf("no references found, this is a root span")
	}
}

func TestTracingRoot(t *testing.T) {
	traceContent := fmt.Sprintf("%x", md5.New().Sum([]byte(time.Now().String())))
	s := startTestServer(nil, 0, func(r *http.Request) {
		th, ok := r.Header[traceHeader]
		if !ok {
			t.Errorf("missing %s request header", traceHeader)
		} else {
			if th[0] != traceContent {
				t.Errorf("wrong X-Trace-Header content: %s", th[0])
			}
		}
	})
	defer s.Close()

	u, _ := url.ParseRequestURI("https://www.example.org/hello")
	r := &http.Request{
		URL:    u,
		Method: "GET",
		Header: make(http.Header),
	}
	w := httptest.NewRecorder()

	doc := fmt.Sprintf(`hello: Path("/hello") -> "%s"`, s.URL)
	tracer := &tracingtest.Tracer{TraceContent: traceContent}
	params := Params{
		OpenTracing: &OpenTracingParams{
			Tracer: tracer,
		},
		Flags: FlagsNone,
	}

	tp, err := newTestProxyWithParams(doc, params)
	if err != nil {
		t.Error(err)
		return
	}
	defer tp.close()

	tp.proxy.ServeHTTP(w, r)

	if len(tracer.RecordedSpans) == 0 {
		t.Fatal("no span recorded...")
	}
	if tracer.RecordedSpans[0].Trace != traceContent {
		t.Errorf("trace not found, got `%s` instead", tracer.RecordedSpans[0].Trace)
	}

	root, ok := tracer.FindSpan("ingress")
	if !ok {
		t.Fatal("root span not found")
	}

	if len(root.Refs) != 0 {
		t.Error("root span cannot have references")
	}
}

func TestTracingSpanName(t *testing.T) {
	traceContent := fmt.Sprintf("%x", md5.New().Sum([]byte(time.Now().String())))
	s := startTestServer(nil, 0, func(r *http.Request) {
		th, ok := r.Header[traceHeader]
		if !ok {
			t.Errorf("missing %s request header", traceHeader)
		} else {
			if th[0] != traceContent {
				t.Errorf("wrong X-Trace-Header content: %s", th[0])
			}
		}
	})
	defer s.Close()

	u, _ := url.ParseRequestURI("https://www.example.org/hello")
	r := &http.Request{
		URL:    u,
		Method: "GET",
		Header: make(http.Header),
	}
	w := httptest.NewRecorder()

	doc := fmt.Sprintf(`hello: Path("/hello") -> tracingSpanName("test-span") -> "%s"`, s.URL)
	tracer := &tracingtest.Tracer{TraceContent: traceContent}
	params := Params{
		OpenTracing: &OpenTracingParams{
			Tracer: tracer,
		},
		Flags: FlagsNone,
	}

	tp, err := newTestProxyWithParams(doc, params)
	if err != nil {
		t.Fatal(err)
	}

	defer tp.close()

	tp.proxy.ServeHTTP(w, r)

	if _, ok := tracer.FindSpan("test-span"); !ok {
		t.Error("setting the span name failed")
	}
}

func TestTracingInitialSpanName(t *testing.T) {
	traceContent := fmt.Sprintf("%x", md5.New().Sum([]byte(time.Now().String())))
	s := startTestServer(nil, 0, func(r *http.Request) {
		th, ok := r.Header[traceHeader]
		if !ok {
			t.Errorf("missing %s request header", traceHeader)
		} else {
			if th[0] != traceContent {
				t.Errorf("wrong X-Trace-Header content: %s", th[0])
			}
		}
	})
	defer s.Close()

	u, _ := url.ParseRequestURI("https://www.example.org/hello")
	r := &http.Request{
		URL:    u,
		Method: "GET",
		Header: make(http.Header),
	}
	w := httptest.NewRecorder()

	doc := fmt.Sprintf(`hello: Path("/hello") -> "%s"`, s.URL)
	tracer := &tracingtest.Tracer{TraceContent: traceContent}
	params := Params{
		OpenTracing: &OpenTracingParams{
			Tracer:      tracer,
			InitialSpan: "test-initial-span",
		},
		Flags: FlagsNone,
	}

	tp, err := newTestProxyWithParams(doc, params)
	if err != nil {
		t.Fatal(err)
	}

	defer tp.close()

	tp.proxy.ServeHTTP(w, r)

	if _, ok := tracer.FindSpan("test-initial-span"); !ok {
		t.Error("setting the span name failed")
	}
}

func TestTracingProxySpan(t *testing.T) {
	const (
		contentSize         = 1 << 16
		prereadSize         = 1 << 12
		responseStreamDelay = 30 * time.Millisecond
	)

	var content bytes.Buffer
	if _, err := io.CopyN(&content, rand.New(rand.NewSource(0)), contentSize); err != nil {
		t.Fatal(err)
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := io.CopyN(w, &content, prereadSize); err != nil {
			t.Fatal(err)
		}

		time.Sleep(responseStreamDelay)
		if _, err := io.Copy(w, &content); err != nil {
			t.Fatal(err)
		}
	}))
	defer s.Close()

	doc := fmt.Sprintf(`* -> "%s"`, s.URL)
	tracer := &tracingtest.Tracer{}
	tp, err := newTestProxyWithParams(doc, Params{OpenTracing: &OpenTracingParams{Tracer: tracer}})
	if err != nil {
		t.Fatal(err)
	}
	defer tp.close()

	req, err := http.NewRequest("GET", "https://www.example.org", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	tp.proxy.ServeHTTP(w, req)

	proxySpan, ok := tracer.FindSpan("proxy")
	if !ok {
		t.Fatal("proxy span not found")
	}

	if proxySpan.FinishTime.Sub(proxySpan.StartTime) < responseStreamDelay {
		t.Error("proxy span did not wait for response stream to finish")
	}
}

func TestTracingProxySpanWithRetry(t *testing.T) {
	const (
		contentSize         = 1 << 16
		prereadSize         = 1 << 12
		responseStreamDelay = 30 * time.Millisecond
	)

	s0 := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	s0.Close()

	content := rand.New(rand.NewSource(0))
	s1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := io.CopyN(w, content, prereadSize); err != nil {
			t.Fatal(err)
		}

		time.Sleep(responseStreamDelay)
		if _, err := io.CopyN(w, content, contentSize-prereadSize); err != nil {
			t.Fatal(err)
		}
	}))
	defer s1.Close()

	const docFmt = `r: * -> <roundRobin, "%s", "%s">;`
	doc := fmt.Sprintf(docFmt, s0.URL, s1.URL)
	tracer := &tracingtest.Tracer{}
	tp, err := newTestProxyWithParams(doc, Params{OpenTracing: &OpenTracingParams{Tracer: tracer}})
	if err != nil {
		t.Fatal(err)
	}
	defer tp.close()

	testFallback := func() bool {
		tracer.Reset("")
		req, err := http.NewRequest("GET", "https://www.example.org", nil)
		if err != nil {
			t.Fatal(err)
		}

		tp.proxy.ServeHTTP(httptest.NewRecorder(), req)

		proxySpans := tracer.FindAllSpans("proxy")
		if len(proxySpans) != 2 {
			t.Log("invalid count of proxy spans", len(proxySpans))
			return false
		}

		for _, s := range proxySpans {
			if s.FinishTime.Sub(s.StartTime) >= responseStreamDelay {
				return true
			}
		}

		t.Log("proxy span with the right duration not found")
		return false
	}

	// Two lb group members are used in round-robin, starting at a non-deterministic index.
	// One of them cannot be connected to, and the proxy should fallback to the other. We
	// want to verify here that the proxy span is traced properly in the fallback case.
	//lint:ignore SA4000 valid testcase in this case
	if !testFallback() && !testFallback() {
		t.Error("failed to trace the right span duration for fallback")
	}
}

func TestProxyTracingDefaultOptions(t *testing.T) {
	t1 := newProxyTracing(nil)
	if t1.tracer == nil || t1.initialOperationName == "" {
		t.Errorf("did not set default options")
	}

	t2 := newProxyTracing(&OpenTracingParams{})
	if t2.tracer == nil || t2.initialOperationName == "" {
		t.Errorf("did not set default options")
	}
}

func TestEnabledLogFilterLifecycleEvents(t *testing.T) {
	tracer := mocktracer.New()
	tracing := newProxyTracing(&OpenTracingParams{
		Tracer:          tracer,
		LogFilterEvents: true,
	})
	span := tracer.StartSpan("test")
	defer span.Finish()

	tracing.logFilterStart(span, "test-filter")
	tracing.logFilterEnd(span, "test-filter")

	mockSpan := span.(*mocktracer.MockSpan)

	if len(mockSpan.Logs()) != 2 {
		t.Errorf("filter lifecycle events were not logged although it was enabled")
	}
}

func TestDisabledLogFilterLifecycleEvents(t *testing.T) {
	tracer := mocktracer.New()
	tracing := newProxyTracing(&OpenTracingParams{
		Tracer:          tracer,
		LogFilterEvents: false,
	})
	span := tracer.StartSpan("test")
	defer span.Finish()

	tracing.logFilterStart(span, "test-filter")
	tracing.logFilterEnd(span, "test-filter")

	mockSpan := span.(*mocktracer.MockSpan)

	if len(mockSpan.Logs()) != 0 {
		t.Errorf("filter lifecycle events were logged although it was disabled")
	}
}
func TestEnabledLogStreamEvents(t *testing.T) {
	tracer := mocktracer.New()
	tracing := newProxyTracing(&OpenTracingParams{
		Tracer:          tracer,
		LogStreamEvents: true,
	})
	span := tracer.StartSpan("test")
	defer span.Finish()

	tracing.logStreamEvent(span, "test-filter", StartEvent)
	tracing.logStreamEvent(span, "test-filter", EndEvent)

	mockSpan := span.(*mocktracer.MockSpan)

	if len(mockSpan.Logs()) != 2 {
		t.Errorf("filter lifecycle events were not logged although it was enabled")
	}
}

func TestDisabledLogStreamEvents(t *testing.T) {
	tracer := mocktracer.New()
	tracing := newProxyTracing(&OpenTracingParams{
		Tracer:          tracer,
		LogStreamEvents: false,
	})
	span := tracer.StartSpan("test")
	defer span.Finish()

	tracing.logStreamEvent(span, "test-filter", StartEvent)
	tracing.logStreamEvent(span, "test-filter", EndEvent)

	mockSpan := span.(*mocktracer.MockSpan)

	if len(mockSpan.Logs()) != 0 {
		t.Errorf("filter lifecycle events were logged although it was disabled")
	}
}

func TestSetEnabledTags(t *testing.T) {
	tracer := mocktracer.New()
	tracing := newProxyTracing(&OpenTracingParams{
		Tracer:      tracer,
		ExcludeTags: []string{},
	})
	span := tracer.StartSpan("test")
	defer span.Finish()

	tracing.setTag(span, HTTPStatusCodeTag, 200)
	tracing.setTag(span, ComponentTag, "skipper")

	mockSpan := span.(*mocktracer.MockSpan)

	tags := mockSpan.Tags()

	_, ok := tags[HTTPStatusCodeTag]
	_, ok2 := tags[ComponentTag]

	if !ok || !ok2 {
		t.Errorf("could not set tags although they were not configured to be excluded")
	}
}

func TestSetDisabledTags(t *testing.T) {
	tracer := mocktracer.New()
	tracing := newProxyTracing(&OpenTracingParams{
		Tracer: tracer,
		ExcludeTags: []string{
			SkipperRouteTag,
		},
	})
	span := tracer.StartSpan("test")
	defer span.Finish()

	tracing.setTag(span, HTTPStatusCodeTag, 200)
	tracing.setTag(span, ComponentTag, "skipper")
	tracing.setTag(span, SkipperRouteTag, "long route definition")

	mockSpan := span.(*mocktracer.MockSpan)

	tags := mockSpan.Tags()

	_, ok := tags[HTTPStatusCodeTag]
	_, ok2 := tags[ComponentTag]
	_, ok3 := tags[SkipperRouteTag]

	if !ok || !ok2 {
		t.Errorf("could not set tags although they were not configured to be excluded")
	}

	if ok3 {
		t.Errorf("a tag was set although it was configured to be excluded")
	}
}

func TestLogEventWithEmptySpan(t *testing.T) {
	tracer := mocktracer.New()
	tracing := newProxyTracing(&OpenTracingParams{
		Tracer: tracer,
	})

	// should not panic
	tracing.logEvent(nil, "test", StartEvent)
	tracing.logEvent(nil, "test", EndEvent)
}

func TestSetTagWithEmptySpan(t *testing.T) {
	tracer := mocktracer.New()
	tracing := newProxyTracing(&OpenTracingParams{
		Tracer: tracer,
	})

	// should not panic
	tracing.setTag(nil, "test", "val")
}
