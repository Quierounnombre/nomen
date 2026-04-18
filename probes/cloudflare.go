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
	record		string
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
		record: "",
	}
	c.obtain_record()
	c.loop()
}

func (c *Cloudflare_probe)loop() {
	for {
		select {
			case cmd := <-c.base.Cmd_ch:
				switch cmd {
				case types.Probe:
					c.execute_probe()
				case types.ShutDown:
					return
				}
		}
	}
}

/*
  DNS resolves?
  No  → DNS-level block (nomen's core case)
  Yes → TCP connects?
          No  → CF down or network issue → hit CF status API
          Yes → HTTP 200?
                  No (403/RST) → ISP proxy block (LaLiga case)
                  Yes → all good
*/

func (c *Cloudflare_probe)execute_probe() {
	var resp types.ProbeResponse
	ok, status := c.base.Http_check()
	if ok {
		switch status {
			case http.StatusOK:
				resp = types.ProbeResponse{ID: c.base.ID, Status: types.StatusOK}
			case http.StatusForbidden:
				resp = types.ProbeResponse{ID: c.base.ID, Status: types.StatusBlocked}
		}
		c.base.Probe_ch <- resp
		return
	}

	if err := c.base.Network_check(); err != nil {
		if err := c.base.Dns_check(); err != nil {
		// DNS broken
		}
		// DNS ok but TCP failed → CF down
	}

	// TCP ok but HTTP not 200 → ISP proxy block (403 etc)
}

func (c *Cloudflare_probe)obtain_record() {
	var result struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	}

	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s", c.region, c.base.Domain),
		nil,
	)
	if err != nil {
		slog.Error("Creating Request", "err", err)
		c.base.Probe_ch <- types.ProbeResponse{ID: c.base.ID, Status: types.StatusError}
		return
	}
	req.Header.Set("Authorization", "Bearer " + c.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Request", "err", err)
		c.base.Probe_ch <- types.ProbeResponse{ID: c.base.ID, Status: types.StatusError}
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		slog.Error("Decoding json", "err", err)
		c.base.Probe_ch <- types.ProbeResponse{ID: c.base.ID, Status: types.StatusError}
		return
	}
	if len(result.Result) > 0 {
		c.record = result.Result[0].ID
		slog.Info("RECORD", "domain", c.record)
	} else {
		slog.Error("No record found", "domain", c.base.Domain)
		c.base.Probe_ch <- types.ProbeResponse{ID: c.base.ID, Status: types.StatusError}
		return
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
