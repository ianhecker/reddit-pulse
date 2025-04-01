package main

import (
	"github.com/vartanbeno/go-reddit/v2/reddit"

	"github.com/ianhecker/reddit-pulse/config"
	"github.com/ianhecker/reddit-pulse/errorChecker"
)

func main() {
	ec := errorChecker.NewErrorChecker()

	cfg, err := config.NewConfig()
	ec.WithMessage("could not make config").CheckErr(err)

	_ := reddit.Credentials{
		ID:       cfg.ClientID,
		Secret:   cfg.ClientSecret,
		Username: cfg.Username,
		Password: cfg.Password,
	}
}
