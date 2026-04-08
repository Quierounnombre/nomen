package probes

import (
	"nomen/types"
	"fmt"
	"net/http"
	"os"
	//"strings"
	"encoding/json"
	"log/slog"
)

type Cloudflare_probe struct {
	base		*types.BaseProbe
	token		string
	region		string
	records		map[string]string
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
		proxy: true,
		records: make(map[string]string),
	}
	c.obtain_records()
	c.loop()
}

func (c *Cloudflare_probe)loop() {
	for {
		select {
			case cmd := <-c.base.Cmd_ch:
				if cmd == types.Probe {
					c.execute_probe()
				}
		}
	}
}

func (c *Cloudflare_probe)execute_probe() {
	var resp types.ProbeResponse

	for _, domain := range c.base.Domains {
		result := basic_probe(domain, c.base.Time_per_probe)
		if result {
			resp = types.ProbeResponse{Name: c.base.Name, Status: types.StatusOK}
		} else {
			resp = types.ProbeResponse{Name: c.base.Name, Status: types.StatusBlocked}
		}
		c.base.Probe_ch <- resp
	}
}

func (c *Cloudflare_probe)obtain_records() {
	var result struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	}

	for _, domain := range c.base.Domains {
		req, err := http.NewRequest("GET",
			fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s", c.region, domain),
			nil,
		)
		req.Header.Set("Authorization", "Bearer " + c.token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
		}
		if len(result.Result) > 0 {
			c.records[domain] = result.Result[0].ID
			slog.Info("RECORD", "domain", c.records[domain])
		}
	}
}

/*
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
*/
