package probes

import (
	"nomen/types"
	"time"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Cloudflare_probe struct {
	base		*types.BaseProbe
	token		string
	region		string
	records		[]string
	proxy		bool
}

func init() {
	types.D_SP["Cloudflare"] = Cloudflare_init
}

func Cloudflare_init(b *types.BaseProbe) {
	c := Cloudflare_probe{
		base: b,
		token: os.Getenv("CF_TOKEN"),
		region: os.Getenv("CF_REGION"),
		records: strings.Split(os.Getenv("CF_RECORD"), ","),
		proxy: true,
	}
	c.loop()
}

func (c *Cloudflare_probe)loop() {
	for {
		select {
			case cmd := <-c.base.Cmd_ch:
				if cmd == types.Probe {
					var resp types.ProbeResponse
					result := basic_probe(c.base.Domain, 5 * time.Minute)
					c.toggle_proxy()
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

func (c *Cloudflare_probe)toggle_proxy() error {
	c.proxy = !c.proxy
	body := fmt.Sprintf(`{"proxied":%v}`, c.proxy)
	req, _ := http.NewRequest("PATCH",
		fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", c.region , c.records[0]),
		strings.NewReader(body),
	)
	req.Header.Set("Authorization", "Bearer " + c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("cloudflare: %d", resp.StatusCode)
	}
	return nil
}
