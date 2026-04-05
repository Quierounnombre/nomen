package main

import (
	"time"
)

type Config struct {
	Provider			[]Provider		`yaml:"provider"`
	Check_interval		time.Duration	`yaml:"check_interval"`
}

type Provider struct {
	Name		string			`yaml:"name"`
	IsPrimary	bool			`yaml:"primary"`
	Domain		string			`yaml:"domain"`
	Status		Status
}

type Status string

const (
	StatusOK		Status = "ok"
	StatusBlocked	Status = "blocked"
	StatusTimeout	Status = "timeout"
)
