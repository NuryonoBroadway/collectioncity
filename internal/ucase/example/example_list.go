// Package example
package example

import (
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/appctx"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/consts"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/repositories"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/ucase/contract"

	"gitlab.privy.id/privypass/privypass-oauth2-core-se/pkg/logger"
)

type exampleList struct {
	repo repositories.Example
}

func NewExampleList(repo repositories.Example) contract.UseCase {
	return &exampleList{repo: repo}
}

// Serve partner list data
func (u *exampleList) Serve(data *appctx.Data) appctx.Response {

	p, err := u.repo.Find(data.Request.Context())

	if err != nil {
		logger.Error(logger.MessageFormat("[example-list] %v", err))

		return *appctx.NewResponse().WithCode(consts.CodeInternalServerError)
	}

	return *appctx.NewResponse().WithCode(consts.CodeSuccess).WithData(p)
}
