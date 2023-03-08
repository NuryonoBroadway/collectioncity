package ucase

import (
	"collection-squad/collection/collection-city/internal/appctx"
	"collection-squad/collection/collection-city/internal/consts"
	"collection-squad/collection/collection-city/internal/ucase/contract"
)

type healthCheck struct {
}

func NewHealthCheck() contract.UseCase {
	return &healthCheck{}
}

func (u *healthCheck) Serve(*appctx.Data) appctx.Response {
	return *appctx.NewResponse().WithCode(consts.CodeSuccess).WithMessage("ok")
}
