package ucase

import (
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/appctx"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/consts"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/ucase/contract"
)

type healthCheck struct {
}

func NewHealthCheck() contract.UseCase {
	return &healthCheck{}
}

func (u *healthCheck) Serve(*appctx.Data) appctx.Response {
	return *appctx.NewResponse().WithCode(consts.CodeSuccess).WithMessage("ok")
}
