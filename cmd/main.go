
package main

import (
	"github.com/rs/zerolog/log"
)

// main executes application
func main() {
	/****** Check wiki for logging best practices ******/
	log.Debug().Msg("Starting up...")
	log.Debug().Msg("Stopping")
}
