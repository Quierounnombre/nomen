package main

import (
	"fmt"
)

func main() {
	set_logger()
	fmt.Println("%+v\n", get_config_from_file_name("config.yaml"))
}
