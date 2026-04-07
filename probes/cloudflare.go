package probes

import (
	"nomen/types"
	"time"
)

type Cloudflare_probe struct {
	base	*types.BaseProbe
}

func init() {
	types.D_SP["Cloudflare"] = Cloudflare_init
}

func Cloudflare_init(b *types.BaseProbe) {
	c := Cloudflare_probe{base: b}
	c.loop()
}

func (c *Cloudflare_probe)loop() {
	for {
		select {
			case cmd := <-c.base.Cmd_ch:
				if cmd == types.Probe {
					var resp types.ProbeResponse
					result := basic_probe(c.base.Domain, 5 * time.Minute)
					if result {
						resp = types.ProbeResponse{Name: c.base.Name, Status: types.StatusOK}
					} else {
						resp = types.ProbeResponse{Name: c.base.Name, Status: types.StatusBlocked}
					}
					c.base.Probe_ch <- resp
				}
		}
	}
}
