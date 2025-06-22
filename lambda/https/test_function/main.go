package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/namikazebadri/app/util"
	"github.com/namikazebadri/lambda/https/test_function/logic"
	"os"
)

func main() {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		lambda.Start(func(ctx context.Context, input json.RawMessage) (string, error) {
			util.Infof("Starting logic from AWS Lambda")

			return logic.LambdaFunction(input)
		})
	} else {
		util.LoadEnv("../../../.env")

		util.Infof("Starting logic from Command Line")

		result, err := logic.LambdaFunction(json.RawMessage{})

		if err != nil {
			util.GetLogger().Panic(err.Error())
		} else {
			util.Infof(result)
		}
	}
}
