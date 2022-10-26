package main

import (
	"admincheckapi/api/config"
	"admincheckapi/api/server"
)

// passed over from linker
var (
	version, build, revision string
)

//
// main loads config, creates the server and starts it
//
func main() {
	config.Init(version, build, revision)
	
	s, err := server.NewAPIServer()
	if err != nil {
		panic(err)
	}

	s.Run()
}
