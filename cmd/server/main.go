package main

import (
	"github.com/ClaudionorJunior/go-expert-api/configs"
)

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
