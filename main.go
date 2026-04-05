package main

import (
	"github.com/goccy/go-yaml"
	"os"
	"fmt"
)

func get_file_content(name string) []byte {
	var content		[]byte
	var err			error

	content, err = os.ReadFile(name)
	if err != nil {
		fmt.Println(err)
	}
	return content
}

func extract_file_content(raw_yaml []byte) *Config {
	var config		Config
	var err			error

	err = yaml.Unmarshal(raw_yaml, &config)
	if err != nil {
		fmt.Println(err)
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

func main() {
	fmt.Printf("%+v\n", get_config_from_file_name("config.yaml"))
}
