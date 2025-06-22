package logic

import (
	"github.com/namikazebadri/app/util"
	"os"
)

type Input struct {
}

func LambdaFunction(input Input) (string, error) {
	util.Infof("Receiving request: ")
	util.Infof("%s", input)

	err := os.Setenv("TZ", util.GetSecretValue("TZ"))

	if err != nil {
		util.Infof("error setting timezone")
	}
	return "Ok", nil
}
