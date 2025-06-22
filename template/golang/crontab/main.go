package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/namikazebadri/app/util"
	"github.com/namikazebadri/template/golang/crontab/logic"
	"os"
)

func main() {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		lambda.Start(func(ctx context.Context) (string, error) {
			util.Infof("Starting logic from AWS Lambda")

			return logic.LambdaFunction(ctx)
		})
	} else {
		util.LoadEnv("../../../.env")

		util.Infof("Starting logic from Command Line")

		ctx := context.TODO()

		result, errLocalCall := logic.LambdaFunction(ctx)

		if errLocalCall != nil {
			util.GetLogger().Panic(errLocalCall.Error())
		} else {
			util.Infof(result)
		}
	}
}
