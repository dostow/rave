package main

import (
	"fmt"

	"github.com/jaffee/commandeer"
	"github.com/dostow/rave/worker"
)

var BUILD string = ""

func main() {
	err := commandeer.Run(worker.NewWorker(BUILD))
	if err != nil {
		fmt.Println(err)
	}
}
