package main

import (
	"nomen/types"
	"nomen/probes"
	"fmt"
	"time"
	"sync"
	"os"
)

func calc_probe_ch_size(config *types.Config) int {
	var l		int

	l = 0
	for _, d := range config.Domains {
		l += len(d.Providers)
	}
	return l
}

func controler(config *types.Config) {
	var probe_response		types.ProbeResponse
	var wg					*sync.WaitGroup

	probe_ch := make(chan types.ProbeResponse, calc_probe_ch_size(config))
	cmds_ch, wg := probes.Init_probes(config, probe_ch)
	ticker := time.Tick(config.Probe_interval)
	for {
		select {
		case probe_response = <-probe_ch:
			switch probe_response.Status {
			case types.StatusOK:
				fmt.Printf("%v\n", probe_response)
			case types.StatusError:
				broadcast(types.ShutDown, cmds_ch)
				wg.Wait()
				os.Exit(1)
			default:
				fmt.Printf("%v\n", probe_response)
			}
		case <-ticker:
			broadcast(types.Probe, cmds_ch)
		}
	}
}

func broadcast(cmd types.Cmd, cmds_ch map[string]chan types.Cmd) {
	for _, ch := range cmds_ch {
		ch <- cmd
	}
}
