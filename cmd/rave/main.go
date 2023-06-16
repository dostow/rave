package main

import (
	"fmt"
	"os"

	"github.com/dostow/rave/handler"
	"github.com/dostow/worker/pkg/worker"
	"github.com/jaffee/commandeer"
)

var BUILD string = ""

func main() {
	command := os.Getenv("COMMAND")
	apiURL := os.Getenv("DOSTOW_API")
	w := worker.NewWorker(&handler.RaveHandler{ApiURL: apiURL})
	w.Name = "push"
	w.Build = BUILD
	w.Command = command

	err := commandeer.Run(w)
	if err != nil {
		fmt.Println(err)
	}
}
