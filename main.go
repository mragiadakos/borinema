package main

import (
	"log"

	"github.com/mragiadakos/borinema/server/conf"
	"github.com/mragiadakos/borinema/server/ctrls"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	println("starting")
	configuration, err := conf.GetConfigurations("conf.toml")
	if err != nil {
		log.Fatal(err)
	}
	ctrls.Run(*configuration)
}
