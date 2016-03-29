package instrumentedhttp

import (
	"net"
	"net/http"
	"sync/atomic"
)

// Stats represents counters about connections
type Stats struct {
	TotalAcceptedConns uint64
	CurrentConns       uint64
}

// ServerInstrumentation provides ConnState callback
// whihch can be attached to any http.Server to collect stats about connections
type ServerInstrumentation struct {
	totalAcceptedConns uint64
	currentConns       uint64
}

// ConnState is a callback for http.Server
func (si *ServerInstrumentation) ConnState(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		atomic.AddUint64(&si.currentConns, 1)
		atomic.AddUint64(&si.totalAcceptedConns, 1)

	case http.StateClosed:
		atomic.AddUint64(&si.currentConns, ^uint64(0))

	case http.StateHijacked:
		atomic.AddUint64(&si.currentConns, ^uint64(0))
	}
}

// Stats return a copy of collected statistics.
func (si *ServerInstrumentation) Stats() Stats {
	return Stats{
		TotalAcceptedConns: atomic.LoadUint64(&si.totalAcceptedConns),
		CurrentConns:       atomic.LoadUint64(&si.currentConns),
	}
}

// ExpvarStats allows to publish Stats as expvar.Func
func (si *ServerInstrumentation) ExpvarStats() interface{} {
	return si.Stats()
}
