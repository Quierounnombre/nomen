package main

import (
	"github.com/goccy/go-yaml"
	"os"
	"log/slog"
)

func get_file_content(name string) []byte {
	var content		[]byte
	var err			error

	content, err = os.ReadFile(name)
	if err != nil {
		slog.Error("failed to read file", "name", name, "err", err)
		os.Exit(1)
	}
	return content
}

func extract_file_content(raw_yaml []byte) *Config {
	var config		Config
	var err			error

	err = yaml.Unmarshal(raw_yaml, &config)
	if err != nil {
		slog.Error("failed to unmarshal content", "raw", raw_yaml, "err", err)
		os.Exit(1)
	}
	return (&config)
}

func get_config_from_file_name(name string) *Config {
	var raw_yaml	[]byte
	var config		*Config

	raw_yaml = get_file_content(name)
	config = extract_file_content(raw_yaml)
	return (config)
}
