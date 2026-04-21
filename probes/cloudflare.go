package probes

import (
	"nomen/types"
	"fmt"
	"net/http"
	"os"
	"strings"
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
				case types.Toggle:
					c.toggle_proxy()
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

	err, status := Http_check(c.base)
	if err != nil {
		err = Network_check(c.base)
		if err != nil {
			err = Dns_check(c.base)
			if err != nil {
				resp = types.ProbeResponse{ID: c.base.ID, Status: types.StatusDnsDown}
				c.base.Probe_ch <- resp
				return
			}
			resp = types.ProbeResponse{ID: c.base.ID, Status: types.StatusNetwork}
			c.base.Probe_ch <- resp
			return
		}
	}
	switch status {
		case http.StatusOK:
			resp = types.ProbeResponse{ID: c.base.ID, Status: types.StatusOK}
		case http.StatusUnauthorized:
			resp = types.ProbeResponse{ID: c.base.ID, Status: types.StatusOK} // Error at auth level, so conectivity is ok
		case http.StatusForbidden:
			resp = types.ProbeResponse{ID: c.base.ID, Status: types.StatusBlocked}
		default:
			slog.Info("STATUS", "esta", status)
			resp = types.ProbeResponse{ID: c.base.ID, Status: types.StatusUnknown}
	}
	c.base.Probe_ch <- resp
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

func (c *Cloudflare_probe)toggle_proxy() {
	c.proxy = !c.proxy
	body := fmt.Sprintf(`{"proxied":%v}`, c.proxy)
	req, err := http.NewRequest("PATCH",
		fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", c.region , c.record),
		strings.NewReader(body),
	)
	if err != nil {
		slog.Error("Creating request", "err", err)
		c.base.Probe_ch <- types.ProbeResponse{ID: c.base.ID, Status: types.StatusError}
		return
	}
	req.Header.Set("Authorization", "Bearer " + c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Problem doing request", "err", err)
		c.base.Probe_ch <- types.ProbeResponse{ID: c.base.ID, Status: types.StatusError}
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		slog.Error("Bad request status", "status", resp.StatusCode)
		c.base.Probe_ch <- types.ProbeResponse{ID: c.base.ID, Status: types.StatusError}
		return
	}
	c.base.Probe_ch <- types.ProbeResponse{ID: c.base.ID, Status: types.StatusOK}
}
