// Package ucase
package ucase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.privy.id/collection/collection-city/internal/appctx"
	"gitlab.privy.id/collection/collection-city/internal/consts"
)

func TestHealthCheck_Serve(t *testing.T) {
	svc := NewHealthCheck()

	t.Run("test health check", func(t *testing.T) {
		result := svc.Serve(&appctx.Data{})

		assert.Equal(t, appctx.Response{
			Code:    consts.CodeSuccess,
			Message: "ok",
		}, result)
	})
}
