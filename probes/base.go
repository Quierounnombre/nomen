package probes

import (
	"log/slog"
	"nomen/types"
	"net/http"
	"time"
	"sync"
	"os"
)

func Init_probes(config *types.Config, probe_ch chan types.ProbeResponse) (map[string]chan types.Cmd, *sync.WaitGroup) {
	var return_ch		map[string]chan types.Cmd
	var wg				sync.WaitGroup

	return_ch = make(map[string]chan types.Cmd)
	for _, domain := range config.Domains {
		for _, p := range domain.Providers {
			_, ok := types.D_SP[p.Name]
			if !ok {
				slog.Error("Provider not supported", "provider", p.Name, "domain", domain.Name)
				os.Exit(1)
			}
			handler := types.D_SP[p.Name]
			cmd_ch := make(chan types.Cmd)
			base_probe := init_base_probe(domain.Name, &domain.Providers[0], probe_ch, cmd_ch)
			return_ch[base_probe.ID] = cmd_ch
			wg.Add(1)
			go func() {
				defer wg.Done()
				handler(base_probe)
			}()
		}
	}
	return return_ch, &wg
}

func init_base_probe(domain string, provider *types.Provider, probe_ch chan types.ProbeResponse, cmd_ch chan types.Cmd) *types.BaseProbe {
	probe := new(types.BaseProbe)
	probe.ID = provider.Name + ":" + domain
	probe.Name = provider.Name
	probe.Status = types.StatusOK
	probe.Current = false
	probe.Cmd_ch = cmd_ch
	probe.Probe_ch = probe_ch
	probe.Capabilities = provider.Capabilities
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
func (b *types.BaseProbe)Dns_check() error {
	r := &net.Resolver{
		PreferGo: true
	}
	ctx, cancel := context.WithTimeout(context.Background(), b.Time_per_probe)
	defer cancel()
	_, err := r.LookupHost(ctx, b.Domain)
	return err
}

//Check if the domain is reachable and accepts connections
func (b *types.BaseProbe)Network_check() error {
	conn, err := net.DialTimeout("tcp", b.Domain+":443", b.Time_per_probe)
	if err == nil {
		conn.Close()
	}
	return err
}

//Check if the service returns a correct status code
func (b *types.BaseProbe)Http_check() (bool, int) {
	client := &http.Client{Timeout: b.Time_per_probe}
	resp, err := client.Get("https://" + b.Domain)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()
	return true, resp.StatusCode
}
