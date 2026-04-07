package main

import (
	"nomen/types"
	"nomen/probes"
	"fmt"
	"time"
)

func controler(config *types.Config) {
	var probe_response		types.ProbeResponse

	cmd_ch := make(chan types.Cmd)
	probe_ch := make(chan types.ProbeResponse)
	probes.Init_probes(config, probe_ch, cmd_ch)
	ticker := time.Tick(config.Probe_interval)
	for {
		select {
		case probe_response = <-probe_ch:
			fmt.Print("%v", probe_response)
		case <-ticker:
			cmd_ch <-types.Probe
		}
	}
}

