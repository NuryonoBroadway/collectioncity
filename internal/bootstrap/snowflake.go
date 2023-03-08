// Package bootstrap
package bootstrap

import (
	"fmt"

	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/helper"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/pkg/generator"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/pkg/logger"
)

// RegistrySnowflake setup snowflake generator
func RegistrySnowflake() {
	hs := helper.GetHostname()
	nodeID := uint64(helper.GetNodeID(hs))

	lf := logger.NewFields(
		logger.EventName("SetupSnowflake"),
		logger.Any("node_id", nodeID),
		logger.Any("hostname", hs),
	)

	logger.Info(fmt.Sprintf(`generate node id for snowflake is %d`, nodeID), lf...)
	generator.Setup(nodeID)
}
