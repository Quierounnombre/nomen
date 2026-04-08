package types

import (
	"time"
)

//-------------------------------------------------------------------------CONFIGS

type Config struct {
	Provider			[]Provider		`yaml:"provider"`
	Probe_interval		time.Duration	`yaml:"probe_interval"`
}

type Provider struct {
	Name			string			`yaml:"name"`
	Capabilities	[]Capability	`yaml:"capabilities"`
	Domains			[]string		`yaml:"domains"`
	Time_per_probe	time.Duration	`yaml:"time_per_probe"`
}

type Capability string

const (
	CapProxyToggle	Capability = "proxy_toggle"
	CapProxyOn		Capability = "proxy_on"
)

//-------------------------------------------------------------------------PROBES

type Status string

const (
	StatusOK		Status = "ok"
	StatusBlocked	Status = "blocked"
	StatusTimeout	Status = "timeout"
	StatusError		Status = "error"
)

type BaseProbe struct {
	Name			string
	Status			Status
	Current			bool
	Cmd_ch			chan Cmd
	Probe_ch		chan ProbeResponse
	Capabilities	[]Capability
	Domains 		[]string
	Time_per_probe	time.Duration
}

type ProbeResponse struct {
	Name			string
	Status			Status
}

//-------------------------------------------------------------------------CMDs

type Cmd string

const (
	ShutDown		Cmd = "shutdown"
	TakeLeadership	Cmd = "leadership"
	Probe			Cmd = "probe"
)
