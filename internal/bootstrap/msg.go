// Package bootstrap
package bootstrap

import (
	"gitlab.privy.id/collection/collection-city/internal/consts"
	"gitlab.privy.id/collection/collection-city/pkg/logger"
	"gitlab.privy.id/collection/collection-city/pkg/msg"
)

func RegistryMessage() {
	err := msg.Setup("msg.yaml", consts.ConfigPath)
	if err != nil {
		logger.Fatal(logger.MessageFormat("file message multi language load error %s", err.Error()))
	}

}
