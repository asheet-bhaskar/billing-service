package worker

import (
	"log"

	"github.com/asheet-bhaskar/billing-service/app/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Start(temporalClient client.Client) {

	w := worker.New(temporalClient, "CREATE_BILL_QUEUE", worker.Options{})

	a := &workflows.Activities{}

	w.RegisterActivity(a.AddLineItemActivity)
	w.RegisterActivity(a.RemoveLineItemActivity)

	w.RegisterWorkflow(workflows.BillingWorkflow)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
