package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/caarlos0/env"
	"github.com/dostow/rave/worker"
	"github.com/jaffee/commandeer"
)

// BUILD contains build id
var BUILD string = "dev"

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetHandler(text.New(os.Stderr))
	w := worker.NewWorker(BUILD)
	if err := env.Parse(w); err != nil {
		log.Debugf("%+v", err)
	}
	err := commandeer.Run(w)
	if err != nil {
		log.Errorf("%+v", err)
	}
}
