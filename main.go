package main

import (
	"github.com/idkarn/curio-db/pkg/server"
)

const PORT uint32 = 3141

func main() {
	config := server.NewConfig(PORT)
	server.Launch(config)
}
