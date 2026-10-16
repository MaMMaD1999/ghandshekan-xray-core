package dispatcher

import (
	"encoding/base64"
	"github.com/xtls/xray-core/common/session"
	"github.com/xtls/xray-core/features/stats"
)

// Extra counters do not rename users, change authentication, or alter routing.
// Base64url makes arbitrary tags/emails unambiguous in the stats namespace.
func scopedTrafficCounter(manager stats.Manager, inbound *session.Inbound, email, direction string, aggregate stats.Counter) stats.Counter {
	if inbound == nil || inbound.Tag == "" || email == "" {
		return aggregate
	}
	encode := base64.RawURLEncoding.EncodeToString
	name := "pronet>>>" + encode([]byte(inbound.Tag)) + ">>>" + encode([]byte(email)) + ">>>traffic>>>" + direction
	detail, err := manager.GetOrRegisterCounter(name)
	if err != nil || detail == nil {
		return aggregate
	}
	return &pairedTrafficCounter{aggregate, detail}
}

type pairedTrafficCounter struct{ aggregate, detail stats.Counter }

func (c *pairedTrafficCounter) Value() int64      { return c.aggregate.Value() }
func (c *pairedTrafficCounter) Set(n int64) int64 { return c.aggregate.Set(n) }
func (c *pairedTrafficCounter) Add(n int64) int64 { c.detail.Add(n); return c.aggregate.Add(n) }
