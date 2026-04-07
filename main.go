package main

import (
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	set_logger()
	godotenv.Load()
	config := get_config_from_file_name("config.yaml")
	fmt.Println("%+v\n", config)
	controler(config)
}
