package machinery

import (
	"testing"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"os"
)

func TestWorker(t *testing.T) {
	type args struct {
		consumerTag string
		workerTasks map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Worker(tt.args.consumerTag, tt.args.workerTasks); (err != nil) != tt.wantErr {
				t.Errorf("Worker() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkerTask(t *testing.T) {
	uri := os.Getenv("redisURI")
	if uri == "" {
		uri = "redis://127.0.0.1:6379"
	}
	println("setting up recoverable machinery with broker " + uri)
	cnf := &config.Config{
		Broker:        uri,
		DefaultQueue:  "machinery_tasks",
		ResultBackend: uri,
	}
	server, err := machinery.NewServer(cnf)
	if err != nil {
		panic(err)
	}
	sig := &tasks.Signature{
		Name: "notification",
		// RetryCount: 3,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: "addonConfig",
			},
			{
				Type:  "string",
				Value: "params",
			},
			{
				Type:  "string",
				Value: "data",
			},
			{
				Type:  "string",
				Value: "trace",
			},
		},
	}
	_, err = server.SendTask(sig)
	// if err != nil {
	// 	t.Errorf("error sending task - %v", err)
	// }

}
