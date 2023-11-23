package main

import (
	"flag"

	"github.com/idkarn/curiodb/pkg/server"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 3141, "Sets the port curiodb will listening on")
	flag.Parse()

	config := server.NewConfig(uint32(port))
	server.Launch(config)
}
