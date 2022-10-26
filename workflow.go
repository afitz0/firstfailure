package firstfailure

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type OrderInfo struct {
	// ...
}

func Workflow(ctx workflow.Context, order OrderInfo) error {
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        time.Second * 120,
		MaximumAttempts:        50,
		NonRetryableErrorTypes: []string{"InsufficientFundsError", "AuthFailure"},
	}

	activityoptions := workflow.ActivityOptions{
		RetryPolicy:         retrypolicy,
		StartToCloseTimeout: 1 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityoptions)

	var a *Activities

	err := workflow.ExecuteActivity(ctx, a.CheckInventory, order).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, a.Charge, order).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, a.FulfillOrder, order).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
