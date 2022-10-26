package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/afitz0/firstfailure"
)

func main() {
	c, err := client.NewLazyClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, firstfailure.TASK_QUEUE, worker.Options{})

	a := &firstfailure.Activities{}
	w.RegisterWorkflow(firstfailure.Workflow)
	w.RegisterActivity(a)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
