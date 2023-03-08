// Package bootstrap
package bootstrap

import (
	"collection-squad/collection/collection-city/internal/consts"
	"collection-squad/collection/collection-city/pkg/logger"
	"collection-squad/collection/collection-city/pkg/msg"
)

func RegistryMessage() {
	err := msg.Setup("msg.yaml", consts.ConfigPath)
	if err != nil {
		logger.Fatal(logger.MessageFormat("file message multi language load error %s", err.Error()))
	}

}
