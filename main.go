package main

import (
	"log"

	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/ctrls"
)

func main() {
	configuration, err := conf.GetConfigurations("conf.toml")
	if err != nil {
		log.Fatal(err)
	}
	ctrls.Run(configuration)
}
