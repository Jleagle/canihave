package logger

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/logging"
	"github.com/Jleagle/canihave/environment"
)

const (
	PROJECT_ID string = "canihaveone-1"
)

func Info(message string) {

	logLocal("Notice: " + message)

	if environment.IsLive() {
		l, client := getLogger()

		l.Log(logging.Entry{
			Severity: logging.Notice,
			Payload:  message,
		})

		go client.Close()
	}
}

func Err(message string) {

	logLocal("Error: " + message)

	if environment.IsLive() {
		l, client := getLogger()

		l.Log(logging.Entry{
			Severity: logging.Error,
			Payload:  message,
		})

		go client.Close()
	}
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

	if environment.IsLocal() {
		return client.Logger("env-local"), client
	}

	return client.Logger("env-live"), client
}

func logLocal(message string) {

	if environment.IsLocal() {
		fmt.Println(message)
	}
}
