package main

import (
	"domofon-api/app"
	"fmt"

	"domofon-api.gg/config"
	"go.uber.org/fx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config       : %v\n", err)
		return
	}

	fx.New(
		fx.Supply(cfg),
		app.App,
	).Run()
}
