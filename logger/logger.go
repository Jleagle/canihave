package logger

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/logging"
)

const (
	PROJECT_ID string = "canihaveone-1"
)

func Info(message string) {

	logLocal("Notice: " + message)

	l, client := getLogger()

	l.Log(logging.Entry{
		Severity: logging.Notice,
		Payload:  message,
	})

	go client.Close()
	//if err != nil {
	//	log.Fatalf("Failed to close client: %v", err)
	//}
}

func Err(message string) {

	logLocal("Error: " + message)

	l, client := getLogger()

	l.Log(logging.Entry{
		Severity: logging.Error,
		Payload:  message,
	})

	go client.Close()
	//if err != nil {
	//	log.Fatalf("Failed to close client: %v", err)
	//}
}

func ErrExit(message string) {
	Err(message)
	os.Exit(1)
}

func getLogger() (logger *logging.Logger, client *logging.Client) {

	ctx := context.Background()
	client, err := logging.NewClient(ctx, PROJECT_ID)
	if err != nil {
		log.Fatalf("Failed to create logging client: %v", err)
	}

	if os.Getenv("CANIHAVE_ENV") == "local" {
		return client.Logger("env-local"), client
	}

	return client.Logger("env-live"), client
}

func logLocal(message string) {

	if os.Getenv("CANIHAVE_ENV") == "local" {
		fmt.Println(message)
	}
}
