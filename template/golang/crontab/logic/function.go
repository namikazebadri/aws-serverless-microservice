package logic

import (
	"context"
	"github.com/namikazebadri/app/util"
	"os"
)

func LambdaFunction(ctx context.Context) (string, error) {
	util.Infof("Receiving request: ")
	util.Infof("%s", ctx)

	err := os.Setenv("TZ", util.GetSecretValue("TZ"))

	if err != nil {
		util.Infof("error setting timezone")
	}

	return "Ok", nil
}
