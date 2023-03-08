package http

import (
	"context"

	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/consts"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/server"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/pkg/logger"
)

// Start function handler starting http listener
func Start(ctx context.Context) {

	serve := server.NewHTTPServer()
	defer serve.Done()
	logger.Info(logger.MessageFormat("starting privypass-oauth2-core-se services... %d", serve.Config().App.Port), logger.EventName(consts.LogEventNameServiceStarting))

	if err := serve.Run(ctx); err != nil {
		logger.Warn(logger.MessageFormat("service stopped, err:%s", err.Error()), logger.EventName(consts.LogEventNameServiceStarting))
	}

	return
}
