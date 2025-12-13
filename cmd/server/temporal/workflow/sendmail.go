package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func SendMailWorkflow(ctx workflow.Context, email string) error {

	logger := workflow.GetLogger(ctx)
	logger.Info("Sending mail to %s", email)
	defer logger.Info("workflow completed")
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 3 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Second * 10,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(ctx, "SendResetPasswordEmail", email).Get(ctx, nil)
	if err != nil {
		logger.Error("Error sending mail", "error", err)
		return err
	}
	return nil

}
