package logger

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/logging"
	"github.com/Jleagle/canihave/environment"
)

const (
	PROJECT_ID string = "canihaveone-1"
)

func Info(message string, err ...error) {
	log(logging.Info, "Notice: "+message, err...)
}

func Err(message string, err ...error) {
	log(logging.Error, "Error: "+message, err...)
}

func ErrExit(message string, err ...error) {
	Err(message, err...)
	os.Exit(1)
}

func log(level logging.Severity, message string, err ...error) {

	if len(err) > 0 {
		message = message + ": " + err[0].Error()
	}

	if environment.IsLive() && false {

		ctx := context.Background()
		c, err := logging.NewClient(ctx, PROJECT_ID)
		if err != nil {
			fmt.Println("Failed to create logging client: " + err.Error())
		}

		c.Logger("all").Log(logging.Entry{
			Severity: level,
			Payload:  message,
		})

		go func() {
			err := c.Close()
			if err != nil {
				fmt.Println("Error sending logs to Google")
			}
		}()

	} else {

		fmt.Println(message)
	}
}
