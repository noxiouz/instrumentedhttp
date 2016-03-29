package instrumentedhttp

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func benchHTTPServer(b *testing.B, useConnState bool) {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	if useConnState {
		var si ServerInstrumentation
		ts.Config.ConnState = si.ConnState
	}
	ts.Start()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		client := &http.Client{}
		for pb.Next() {
			resp, err := client.Get(ts.URL)
			if err != nil {
				b.Errorf("%v", err)
			}
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}
	})
}

func BenchmarkVanilaHTTPServer(b *testing.B) {
	benchHTTPServer(b, false)
}

func BenchmarkInstrumentedHTTPServer(b *testing.B) {
	benchHTTPServer(b, true)
}
