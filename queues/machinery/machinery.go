package machinery

// TODO: move this to separate go mod
import (
	// "context"
	// "errors"
	// "fmt"
	"os"
	// "time"

	// opentracing "github.com/opentracing/opentracing-go"
	// opentracing_log "github.com/opentracing/opentracing-go/log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	// "github.com/google/uuid"
)

// DefaultServer default machinery server
var DefaultServer *machinery.Server

func loadConfig() (*config.Config, error) {
	// if configPath != "" {
	// 	return config.NewFromYaml(configPath, true)
	// }

	return config.NewFromEnvironment(true)
}

func makeServer() (*machinery.Server, error) {
	// cnf, err := loadConfig()
	// if err != nil {
	// 	return nil, err
	// }

	// // Create server instance
	// server, err := machinery.NewServer(cnf)
	// if err != nil {
	// 	return nil, err
	// }
	// return server, nil
	uri := os.Getenv("redisURI")
	if uri == "" {
		uri = "redis://127.0.0.1:6379"
	}
	println("setting broker " + uri)
	cnf := &config.Config{
		Broker:        uri,
		DefaultQueue:  "machinery_tasks",
		ResultBackend: uri,
	}
	return machinery.NewServer(cnf)
}
func startServer(tasks map[string]interface{}) (*machinery.Server, error) {
	server, err := makeServer()
	if err != nil {
		return nil, err
	}
	return server, server.RegisterTasks(tasks)
}

// Worker a machinery worker
func Worker(consumerTag string, workerTasks map[string]interface{}) error {
	// cleanup, err := tracers.SetupTracer(consumerTag)
	// if err != nil {
	// 	log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	// }
	// defer cleanup()

	server, err := startServer(workerTasks)
	if err != nil {
		return err
	}
	DefaultServer = server

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(consumerTag, 0)

	// Here we inject some custom code for error handling,
	// start and end of task hooks, useful for metrics for example.
	errorhandler := func(err error) {
		log.ERROR.Println("error:", err)
	}

	pretaskhandler := func(signature *tasks.Signature) {
		log.INFO.Println("start:", signature.Name)
	}

	posttaskhandler := func(signature *tasks.Signature) {
		log.INFO.Println("end:", signature.Name)
	}

	worker.SetPostTaskHandler(posttaskhandler)
	worker.SetErrorHandler(errorhandler)
	worker.SetPreTaskHandler(pretaskhandler)
	return worker.Launch()
}
