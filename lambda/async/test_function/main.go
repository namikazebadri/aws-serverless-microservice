package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/namikazebadri/app/util"
	"github.com/namikazebadri/lambda/async/test_function/logic"
	"os"
)

func main() {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		lambda.Start(func(ctx context.Context, event events.SQSEvent) (string, error) {
			util.Infof("Starting logic from AWS Lambda")

			var input logic.Input

			for _, message := range event.Records {
				util.Infof("Message ID: %s", message.MessageId)
				util.Infof("Message Body: %s", message.Body)

				err := json.Unmarshal([]byte(message.Body), &input)

				if err != nil {
					return "", err
				}
			}

			return logic.LambdaFunction(input)
		})
	} else {
		util.LoadEnv("../../../.env")

		util.Infof("Starting logic from Command Line")

		result, errCall := logic.LambdaFunction(logic.Input{})

		if errCall != nil {
			util.GetLogger().Panic(errCall.Error())
		} else {
			util.Infof(result)
		}
	}
}
