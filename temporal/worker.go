package main

import (
	"concert/internal/utils"
	"concert/temporal/activity"
	"concert/temporal/workflow"
	"fmt"
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
		return
	}
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	if smtpHost == "" || smtpPort == "" {
		fmt.Printf("could not configure the smpt server %s:%", smtpHost, smtpPort)
		return
	}
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
