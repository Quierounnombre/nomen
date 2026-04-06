package types

import (
	"time"
)

//-------------------------------------------------------------------------CONFIGS

type Config struct {
	Provider			[]Provider		`yaml:"provider"`
	Domain				string			`yaml:"domain"`
	Check_interval		time.Duration	`yaml:"check_interval"`
}

type Provider struct {
	Name			string			`yaml:"name"`
	Capabilities	[]Capability	`yaml:"capabilities"`
}

type Capability string

const (
	CapProxyToggle	Capability = "proxy_toggle"
)

//-------------------------------------------------------------------------PROBES

type Status string

const (
	StatusOK		Status = "ok"
	StatusBlocked	Status = "blocked"
	StatusTimeout	Status = "timeout"
)

type BaseProbe struct {
	Name			string
	Status			Status
	Current			bool
	Cmd_ch			chan Cmd
	Probe_ch		chan ProbeResponse
	Capabilities	[]Capability
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
)
