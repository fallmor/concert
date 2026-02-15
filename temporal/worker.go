package main

import (
	"concert/internal/utils"
	"concert/temporal/activity"
	"concert/temporal/workflow"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	projectRoot := utils.GetProjectRoot("go.mod")
	envPath := filepath.Join(projectRoot, ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	temporalHost := utils.GetEnvOrDefault("TEMPORAL_HOST", "localhost:7233")
	c, err := client.Dial(client.Options{
		HostPort: temporalHost,
	})
	if err != nil {
		log.Fatalf("failed to connect to temporal server: %v", err)
	}
	defer c.Close()

	resendAPI := os.Getenv("RESEND_API")
	if resendAPI == "" {
		log.Fatal("RESEND_API environment variable is required")
	}
	sendmail := activity.NewEmailActivities(resendAPI)

	w := worker.New(c, "email-task-queue", worker.Options{})
	w.RegisterWorkflow(workflow.SendMailWorkflow)
	w.RegisterActivity(sendmail.SendResetPasswordEmail)
	w.Run(worker.InterruptCh())
}
