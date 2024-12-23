package workflows

// import (
// 	"context"
// 	"log"

// 	"github.com/mitchellh/mapstructure"
// 	"go.temporal.io/sdk/workflow"
// )

// type AddBillItemSignal struct {
// 	BillID   int64
// 	BillItem domain.BillItem
// }

// type RemoveBillItemSignal struct {
// 	BillItemID int64
// 	BillID     int64
// }

// type CloseBillSignal struct {
// 	BillID int64
// }

// func createBillActivity(_ context.Context, bill domain.Bill) error {
// 	return nil
// }

// func addBillItemActivity(_ context.Context, billID int64, billItem domain.BillItem) error {

// 	return nil
// }

// func removeBillItemActivity(_ context.Context, billItemID int64) error {
// 	log.Println("removing bill item")
// 	return nil
// }

// func closeBillActivity(_ context.Context, billID int64) error {
// 	log.Println("closing bill")
// 	return nil
// }

// func BillingWorkflow(ctx workflow.Context, bill *domain.Bill) error {
// 	logger := workflow.GetLogger(ctx)

// 	addBillItemChan := workflow.GetSignalChannel(ctx, "ADD_BILL_ITEM_CHANNEL")
// 	removeBillItemChan := workflow.GetSignalChannel(ctx, "REMOVE_BILL_ITEM_CHANNEL")
// 	CloseBillChan := workflow.GetSignalChannel(ctx, "CLOSE_BILL_CHANNEL")

// 	for {
// 		selector := workflow.NewSelector(ctx)

// 		selector.AddReceive(addBillItemChan, func(c workflow.ReceiveChannel, _ bool) {
// 			var signal interface{}
// 			c.Receive(ctx, &signal)

// 			var message AddBillItemSignal
// 			err := mapstructure.Decode(signal, &message)
// 			if err != nil {
// 				logger.Error("Invalid signal type %v", err)
// 				return
// 			}

// 			ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{})

// 			err = workflow.ExecuteActivity(ctx, addBillItemActivity, message.BillID, message.BillItem).Get(ctx, nil)
// 			if err != nil {
// 				logger.Error("Error adding bill item: %v", err)
// 				return
// 			}
// 		})

// 		selector.AddReceive(removeBillItemChan, func(c workflow.ReceiveChannel, _ bool) {
// 			var signal interface{}
// 			c.Receive(ctx, &signal)

// 			var message RemoveBillItemSignal
// 			err := mapstructure.Decode(signal, &message)
// 			if err != nil {
// 				logger.Error("Invalid signal type %v", err)
// 				return
// 			}

// 			ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{})
// 			err = workflow.ExecuteActivity(ctx, removeBillItemActivity, message).Get(ctx, nil)
// 			if err != nil {
// 				logger.Error("Error removing bill item: %v", err)
// 				return
// 			}
// 		})

// 		selector.AddReceive(CloseBillChan, func(c workflow.ReceiveChannel, _ bool) {
// 			var signal interface{}
// 			c.Receive(ctx, &signal)

// 			var message CloseBillSignal
// 			err := mapstructure.Decode(signal, &message)
// 			if err != nil {
// 				logger.Error("Invalid signal type %v", err)
// 				return
// 			}

// 			ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{})

// 			err = workflow.ExecuteActivity(ctx, closeBillActivity, message).Get(ctx, nil)
// 			if err != nil {
// 				logger.Error("Error closing bill: %v", err)
// 				return
// 			}
// 		})

// 		selector.Select(ctx)
// 	}
// }
