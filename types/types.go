package types

import (
	"time"
)

//-------------------------------------------------------------------------CONFIGS

type Config struct {
	Domains				[]Domains		`yaml:"domains"`
	Probe_interval		time.Duration	`yaml:"probe_interval"`
}

type Domains struct {
	Providers			[]Provider		`yaml:"providers"`
	Name				string			`yaml:"name"`
}

type Provider struct {
	Name			string			`yaml:"name"`
	Capabilities	[]Capability	`yaml:"capabilities"`
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
	StatusUnknown	Status = "Unknown"
	StatusDnsDown	Status = "dns"
	StatusNetwork	Status = "network"
)

type BaseProbe struct {
	ID				string
	Name			string
	Cmd_ch			chan Cmd
	Probe_ch		chan ProbeResponse
	Domain			string
	Time_per_probe	time.Duration
}

type ProbeResponse struct {
	ID				string
	Status			Status
}

type ProviderState struct {
	Status			Status
	Capabilities	[]Capability
	Cmd_ch			chan Cmd
}

type DomainState struct {
	Current			string
	Providers		map[string]ProviderState
}

//-------------------------------------------------------------------------CMDs

type Cmd string

const (
	ShutDown		Cmd = "shutdown"
	TakeLeadership	Cmd = "leadership"
	Probe			Cmd = "probe"
	Toggle			Cmd = "toggle"
)
