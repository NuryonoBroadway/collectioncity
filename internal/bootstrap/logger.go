// Package bootstrap
package bootstrap

import (
	"gitlab.privy.id/collection/collection-city/internal/appctx"
	"gitlab.privy.id/collection/collection-city/pkg/logger"
	"gitlab.privy.id/collection/collection-city/pkg/util"
)

func RegistryLogger(cfg *appctx.Config) {
	logger.Setup(logger.Config{
		Environment: util.EnvironmentTransform(cfg.App.Env),
		Debug:       cfg.App.Debug,
		Level:       cfg.Logger.Level,
		ServiceName: cfg.Logger.Name,
	})
}
