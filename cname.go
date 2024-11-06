package cname

import (
	"context"
	"math"

	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

// Rewrite is plugin to rewrite requests internally before being handled.
type Cname struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (c Cname) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	mw := NewResponseModifier(w)
	return plugin.NextOrFailure(c.Name(), c.Next, ctx, mw, r)
}

// Name implements the Handler interface.
func (c Cname) Name() string { return "cname" }

type ResponseModifier struct {
	dns.ResponseWriter
}

// Returns a dns.Msg modifier that replaces CNAME on root zones with other records.
func NewResponseModifier(w dns.ResponseWriter) *ResponseModifier {
	return &ResponseModifier{
		ResponseWriter: w,
	}
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// WriteMsg records the status code and calls the
// underlying ResponseWriter's WriteMsg method.
func (r *ResponseModifier) WriteMsg(res *dns.Msg) error {
	// Find and delete CNAME record on that zone, storing the canonical name.
	var (
		cname string
		ttl   uint32 = math.MaxUint32
	)

	for i := 0; i < len(res.Answer); {
		rr := res.Answer[i]
		if rr.Header().Rrtype == dns.TypeCNAME {
			cname = rr.(*dns.CNAME).Hdr.Name
			ttl = min(ttl, rr.(*dns.CNAME).Header().Ttl)
			// Remove the CNAME record
			res.Answer = append(res.Answer[:i], res.Answer[i+1:]...)
			continue
		}
		i++
	}

	// Rename all the records with the above canonical name to the zone name
	for _, rr := range res.Answer {
		if cname != "" {
			rr.Header().Name = cname
			rr.Header().Ttl = min(ttl, rr.Header().Ttl)
		}
	}

	return r.ResponseWriter.WriteMsg(res)
}

// Write is a wrapper that records the size of the message that gets written.
func (r *ResponseModifier) Write(buf []byte) (int, error) {
	n, err := r.ResponseWriter.Write(buf)
	return n, err
}

// Hijack implements dns.Hijacker. It simply wraps the underlying
// ResponseWriter's Hijack method if there is one, or returns an error.
func (r *ResponseModifier) Hijack() {
	r.ResponseWriter.Hijack()
}
