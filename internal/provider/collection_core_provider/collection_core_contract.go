package collection_core_provider

import (
	"context"

	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/entity"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/provider"
	"gitlab.privy.id/privypass/privypass-package-core/request/meta"
)

type CollectionCoreProvider interface {
	CreateCollectionCityGrpc(ctx context.Context, city entity.City) error
	CreateCollectionCityPubsub(ctx context.Context, city entity.City) error
	GetCollectionCityAll(ctx context.Context, m meta.Metadata) ([]entity.City, *meta.Metadata, error)
	GetCollectionCityByDocumentID(ctx context.Context, documentID string) (*entity.City, error)
}

type CollectionCore struct {
	Provider *provider.Provider
}

func NewCollecionCoreProvider(Provider *provider.Provider) CollectionCoreProvider {
	return &CollectionCore{Provider: Provider}
}
