package workflows

import (
	"time"

	"github.com/asheet-bhaskar/billing-service/app/models"
	"github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/workflow"
)

type LineItemSignal struct {
	BillID string
	ItemID string
}

func BillingWorkflow(ctx workflow.Context, bill *models.Bill) error {
	logger := workflow.GetLogger(ctx)

	var a *Activities
	addLineItemChan := workflow.GetSignalChannel(ctx, "ADD_BILL_ITEM_CHANNEL")
	removeLineItemChan := workflow.GetSignalChannel(ctx, "REMOVE_BILL_ITEM_CHANNEL")

	for {
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(addLineItemChan, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			var message LineItemSignal
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			ao := workflow.ActivityOptions{
				StartToCloseTimeout: time.Minute,
			}
			ctx = workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(ctx, a.AddLineItemActivity, message).Get(ctx, nil)
			if err != nil {
				logger.Error("Error adding bill item: %v", err)
				return
			}
		})

		selector.AddReceive(removeLineItemChan, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			var message LineItemSignal
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			ao := workflow.ActivityOptions{
				StartToCloseTimeout: time.Minute,
			}
			ctx = workflow.WithActivityOptions(ctx, ao)
			err = workflow.ExecuteActivity(ctx, a.RemoveLineItemActivity, message).Get(ctx, nil)
			if err != nil {
				logger.Error("Error removing bill item: %v", err)
				return
			}
		})

		selector.Select(ctx)
	}
}
