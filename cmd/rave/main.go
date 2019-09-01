package main

import (
	"fmt"

	"github.com/apex/log"
	"os"
	"github.com/apex/log/handlers/text"
	"github.com/jaffee/commandeer"
	"github.com/dostow/rave/worker"
)

var BUILD string = ""

func main() {  
	log.SetLevel(log.DebugLevel)
	log.SetHandler(text.New(os.Stderr)) 
	err := commandeer.Run(worker.NewWorker(BUILD))
	if err != nil {
		fmt.Println(err)
	}
}
