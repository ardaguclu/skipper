/*
This command provides an executable version of skipper with the default
set of filters.

For the list of command line options, run:

    skipper -help

For details about the usage and extensibility of skipper, please see the
documentation of the root skipper package.

To see which built-in filters are available, see the skipper/filters
package documentation.
*/
package main

import (
	"fmt"

	"github.com/ardaguclu/skipper"
	"github.com/ardaguclu/skipper/config"
	log "github.com/sirupsen/logrus"
)

var (
	version string
	commit  string
)

func main() {
	cfg := config.NewConfig()
	if err := cfg.Parse(); err != nil {
		log.Fatalf("Error processing config: %s", err)
	}

	if cfg.PrintVersion {
		fmt.Printf(
			"Skipper version %s (commit: %s)\n",
			version, commit,
		)

		return
	}

	log.SetLevel(cfg.ApplicationLogLevel)
	log.Fatal(skipper.Run(cfg.ToOptions()))
}
