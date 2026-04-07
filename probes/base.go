package probes

import (
	"log/slog"
	"nomen/types"
	"net/http"
	"time"
)

func Init_probes(config *types.Config, probe_ch chan types.ProbeResponse, cmd_ch chan types.Cmd) {
	for i, _ := range config.Provider {
		provider := config.Provider[i]
		handler, ok := types.D_SP[provider.Name]
		if !ok {
			slog.Error("Provider not supported(yet), PR welcome!") 
			continue
		}
		base_probe := init_base_probe(&provider, probe_ch, cmd_ch, config.Domain)
		go handler(base_probe)
	}
}

func init_base_probe(provider *types.Provider, probe_ch chan types.ProbeResponse, cmd_ch chan types.Cmd, domain string) *types.BaseProbe {
	probe := new(types.BaseProbe)
	probe.Name = provider.Name
	probe.Status = types.StatusOK
	probe.Current = false
	probe.Cmd_ch = cmd_ch
	probe.Probe_ch = probe_ch
	probe.Capabilities = provider.Capabilities
	probe.Domain = domain
	return (probe)
}

//Basic probes that checks if a domain is reachable
func basic_probe(domain string, timeout time.Duration) bool {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get("https://" + domain)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}
