package temporal

import (
	"concert/cmd/server/temporal/activity"
	"concert/cmd/server/temporal/workflow"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_Port")
	smtpUser := os.Getenv("SMTP_User")
	smtpPass := os.Getenv("SMTP_Pass")
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalf("failed to connect to temporal server: %v", err)
	}
	defer c.Close()
	sendmail := activity.NewEmailActivities(smtpHost, smtpPort, smtpUser, smtpPass)

	w := worker.New(c, "email-task-queue", worker.Options{})
	w.RegisterWorkflow(workflow.SendMailWorkflow)
	w.RegisterActivity(sendmail.SendResetPasswordEmail)
	w.Run(worker.InterruptCh())
}
