package main

import (
	"nomen/types"
	"nomen/probes"
	"fmt"
	"time"
	"sync"
	"os"
)

func controler(config *types.Config) {
	var probe_response		types.ProbeResponse
	var wg					*sync.WaitGroup

	probe_ch := make(chan types.ProbeResponse)
	cmds_ch, wg := probes.Init_probes(config, probe_ch)
	ticker := time.Tick(config.Probe_interval)
	for {
		select {
		case probe_response = <-probe_ch:
			switch probe_response.Status {
			case types.StatusOK:
				fmt.Printf("%v", probe_response)
			case types.StatusError:
				broadcast(types.ShutDown, cmds_ch)
				wg.Wait()
				os.Exit(1)
			}
		case <-ticker:
			broadcast(types.Probe, cmds_ch)
		}
	}
}

func broadcast(cmd types.Cmd, cmds_ch []chan types.Cmd) {
	for _, ch := range cmds_ch {
		ch <- cmd
	}
}
