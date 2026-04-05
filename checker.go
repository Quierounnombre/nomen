package main

import (
	"log/slog"
)

func init_checkers(config *Config) {
	for i, _ := range config {
		go start_provider(c, i)
	}
}

func start_provider(config *Config, index int) {
	provider := config.Provider[i]
	handler, ok := D_SP[provider.Name]
	if !ok {
		slog.Error("Provider not supported(yet), PR welcome!") 
		return
	}
	handler(provider)
}
