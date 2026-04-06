package probes

import (
	"log/slog"
	"nomen/types"
)

func Init_probes(config *types.Config, probe_ch chan types.ProbeResponse, cmd_ch chan types.Cmd) {
	for i, _ := range config.Provider {
		provider := config.Provider[i]
		handler, ok := types.D_SP[provider.Name]
		if !ok {
			slog.Error("Provider not supported(yet), PR welcome!") 
			continue
		}
		base_probe := init_base_probe(&provider, probe_ch, cmd_ch)
		go handler(base_probe)
	}
}

func init_base_probe(provider *types.Provider, probe_ch chan types.ProbeResponse, cmd_ch chan types.Cmd) *types.BaseProbe {
	probe := new(types.BaseProbe)
	probe.Name = provider.Name
	probe.Status = types.StatusOK
	probe.Current = false
	probe.Cmd_ch = cmd_ch
	probe.Probe_ch = probe_ch
	probe.Capabilities = provider.Capabilities
	return (probe)
}
