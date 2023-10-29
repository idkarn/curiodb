package main

import (
	"github.com/idkarn/curiodb/pkg/server"
)

const PORT uint32 = 3141

func main() {
	config := server.NewConfig(PORT)
	server.Launch(config)
}
