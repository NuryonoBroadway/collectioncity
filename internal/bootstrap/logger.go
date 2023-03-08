// Package bootstrap
package bootstrap

import (
	"collection-squad/collection/collection-city/internal/appctx"
	"collection-squad/collection/collection-city/pkg/logger"
	"collection-squad/collection/collection-city/pkg/util"
)

func RegistryLogger(cfg *appctx.Config) {
	logger.Setup(logger.Config{
		Environment: util.EnvironmentTransform(cfg.App.Env),
		Debug:       cfg.App.Debug,
		Level:       cfg.Logger.Level,
		ServiceName: cfg.Logger.Name,
	})
}
