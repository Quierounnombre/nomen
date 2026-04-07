package main

import (
	"fmt"
)

func main() {
	set_logger()
	config := get_config_from_file_name("config.yaml")
	fmt.Println("%+v\n", config)
	controler(config)
}
