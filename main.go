package main

import (
	"referral-service/app"
	"referral-service/controller"
	"referral-service/handler"
	"referral-service/repository"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.Module,        // provide gateways.
		repository.Module, // provide reposity interface.
		controller.Module, // provide controller interface.
		handler.Module,    // wire up to handlers.
	).Run()
}
