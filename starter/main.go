package main

import (
	"context"
	"fmt"
	"log"

	"go.temporal.io/sdk/client"

	"github.com/afitz0/firstfailure"
)

const workflows = 100

func main() {
	c, err := client.NewLazyClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	for i := 0; i < workflows; i++ {
		workflowOptions := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("failure-run-%d", i),
			TaskQueue: firstfailure.TASK_QUEUE,
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, firstfailure.Workflow, firstfailure.OrderInfo{})
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}

		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	}

}
