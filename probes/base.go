package probes

import (
	"log/slog"
	"nomen/types"
	"net/http"
	"net"
	"time"
	"sync"
	"os"
	"context"
)

//Need to return a map of domainsstates each one initialize with the data and probes
func Init_probes(config *types.Config, probe_ch chan types.ProbeResponse) (map[string]types.DomainState, *sync.WaitGroup) {
	var return_ds	map[string]types.DomainState
	var wg			sync.WaitGroup

	return_ds = make(map[string]types.DomainState)
	for _, domain := range config.Domains {
		return_ds[domain.Name] = types.DomainState {
			Providers: make(map[string]types.ProviderState),
		}
		for _, p := range domain.Providers {
			_, ok := types.D_SP[p.Name]
			if !ok {
				slog.Error("Provider not supported", "provider", p.Name, "domain", domain.Name)
				os.Exit(1)
			}
			handler := types.D_SP[p.Name]
			cmd_ch := make(chan types.Cmd)
			base_probe := init_base_probe(domain.Name, &p, probe_ch, cmd_ch)
			ds := return_ds[domain.Name]
			ds.Providers[p.Name] = types.ProviderState {
				Cmd_ch: cmd_ch,
				Capabilities: p.Capabilities,
				Status: types.StatusUnknown,
			}
			if return_ds[domain.Name].Current == "" {
				ds.Current = p.Name
			}
			return_ds[domain.Name] = ds
			wg.Add(1)
			go func() {
				defer wg.Done()
				handler(base_probe)
			}()
		}
	}
	return return_ds, &wg
}

func init_base_probe(domain string, provider *types.Provider, probe_ch chan types.ProbeResponse, cmd_ch chan types.Cmd) *types.BaseProbe {
	probe := new(types.BaseProbe)
	probe.ID = provider.Name + ":" + domain
	probe.Name = provider.Name
	probe.Cmd_ch = cmd_ch
	probe.Probe_ch = probe_ch
	probe.Domain = domain
	probe.Time_per_probe = provider.Time_per_probe
	return (probe)
}

//Basic probes that checks if a domain is reachable, just dont check for any posible error be aware
func basic_probe(domain string, timeout time.Duration) bool {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get("https://" + domain)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

//Check if the DNS resolve, only the DNS
func Dns_check(b *types.BaseProbe) error {
	r := &net.Resolver{
		PreferGo: true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.Time_per_probe)
	defer cancel()
	_, err := r.LookupHost(ctx, b.Domain)
	return err
}

//Check if the domain is reachable and accepts connections
func Network_check(b *types.BaseProbe) error {
	conn, err := net.DialTimeout("tcp", b.Domain+":443", b.Time_per_probe)
	if err == nil {
		conn.Close()
	}
	return err
}

//Check if the service returns a correct status code
func Http_check(b *types.BaseProbe) (error, int) {
	client := &http.Client{Timeout: b.Time_per_probe}
	resp, err := client.Get("https://" + b.Domain)
	if err != nil {
		return err, 0
	}
	defer resp.Body.Close()
	return err, resp.StatusCode
}
