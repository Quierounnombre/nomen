package main

import (
	"nomen/types"
	"nomen/probes"
	"strings"
	"fmt"
	"time"
	"sync"
	"os"
	"log/slog"
)

func calc_probe_ch_size(config *types.Config) int {
	var l		int

	l = 0
	for _, d := range config.Domains {
		l += len(d.Providers)
	}
	return l
}

func select_next_healthy(domain types.DomainState) string {
	for provider range domain.Providers {
		if domain[provider].Status == StatusOK {
			return provider
		}
	}
	return ""
}

func controler(config *types.Config) {
	var probe_response		types.ProbeResponse
	var wg					*sync.WaitGroup

	probe_ch := make(chan types.ProbeResponse, calc_probe_ch_size(config))
	ds, wg := probes.Init_probes(config, probe_ch)
	ticker := time.Tick(config.Probe_interval)
	for {
		select {
		case probe_response = <-probe_ch:
			provider, domain, _ := strings.Cut(probe_response.ID, ":")
			state := ds[domain]
			switch probe_response.Status {
				case types.StatusOK:
					slog.Info("OK", "ID", probe_response.ID)
				case types.StatusError:
					broadcast(types.ShutDown, ds)
					wg.Wait()
					os.Exit(1)
				case types.StatusBlocked:
					next := select_next_healthy(state)
					if state.Current == "Cloudflare" {
					}
					if provider == "Cloudflare" {
						send(types.Toggle, ds, probe_response.ID)
					}
				default:
					fmt.Printf("%v\n", probe_response)
			}
		case <-ticker:
			broadcast(types.Probe, ds)
		}
	}
}

func send(cmd types.Cmd, ds map[string]types.DomainState, ID string) {
	domain, provider, _ := strings.Cut(ID, ":")
	ds[domain].Providers[provider].Cmd_ch <- cmd
}

func broadcast(cmd types.Cmd, ds map[string]types.DomainState) {
	for _, domain := range ds {
		for _, provider := range domain.Providers {
			provider.Cmd_ch <- cmd
		}
	}
}
