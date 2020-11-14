package main

import (
	// system

	"github.com/aufheben/mutuals-server/local/routing"
)

func main() {
	routing.Init()
	routing.Launch()
}
