package main

import (
	"time"
)

type Config struct {
	Provider			[]Provider		`yaml:"provider"`
	Domain				string			`yaml:"domain"`
	Check_interval		time.Duration	`yaml:"check_interval"`
}

type Provider struct {
	Name			string			`yaml:"name"`
	IsPrimary		bool			`yaml:"primary"`
	Capabilities	[]Capability	`yaml:"capabilities"`
	Status			Status
	Current			bool
}

type Capability string

const (
	CapProxyToggle	Capability = "proxy_toggle"
)

type Status string

const (
	StatusOK		Status = "ok"
	StatusBlocked	Status = "blocked"
	StatusTimeout	Status = "timeout"
)
