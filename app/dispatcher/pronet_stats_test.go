package dispatcher

import (
	"github.com/xtls/xray-core/common/session"
	"testing"
)

type testCounter struct{ n int64 }

func (c *testCounter) Value() int64      { return c.n }
func (c *testCounter) Set(v int64) int64 { old := c.n; c.n = v; return old }
func (c *testCounter) Add(v int64) int64 { c.n += v; return c.n }
func TestPairedAccountingDoesNotDuplicateAggregate(t *testing.T) {
	total, a, b := &testCounter{}, &testCounter{}, &testCounter{}
	(&pairedTrafficCounter{total, a}).Add(100)
	(&pairedTrafficCounter{total, b}).Add(40)
	if total.Value() != 140 || a.Value() != 100 || b.Value() != 40 {
		t.Fatal(total, a, b)
	}
	(&pairedTrafficCounter{total, a}).Set(0)
	if a.Value() != 100 {
		t.Fatal("aggregate reset erased independent inbound counter")
	}
}
func TestMissingScopeLeavesCounterUnchanged(t *testing.T) {
	c := &testCounter{}
	if scopedTrafficCounter(nil, &session.Inbound{}, "u", "uplink", c) != c {
		t.Fatal("empty tag")
	}
}
