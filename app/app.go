package app

import (
	"go.uber.org/fx"

	"referral-service/config"
	"referral-service/logger"
)

var Module = fx.Module(
	"app base",
	logger.Module,
	config.Module,
)
