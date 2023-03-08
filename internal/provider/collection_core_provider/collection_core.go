package collection_core_provider

import (
	"context"
	"fmt"

	"collection-squad/collection/collection-city/internal/entity"

	"gitlab.privy.id/privypass/privypass-package-core/collection_core"
	"gitlab.privy.id/privypass/privypass-package-core/request/meta"
)

func (prc *CollectionCore) CreateCollectionCityGrpc(ctx context.Context, city entity.City) error {
	var query = prc.Provider.CollectionCore.NewQuery().Col("cities").Doc(city.ID)
	defer query.Close()

	_, err := query.Set(ctx, city, collection_core.IsMergeAll, collection_core.IsUseGRPC)
	if err != nil {
		return err
	}

	return nil
}

func (prc *CollectionCore) CreateCollectionCityPubsub(ctx context.Context, city entity.City) error {
	var query = prc.Provider.CollectionCore.NewQuery().Col("cities").Doc(city.ID)
	defer query.Close()

	_, err := query.Set(ctx, city, collection_core.IsMergeAll, collection_core.IsUsePubSub)
	if err != nil {
		return err
	}

	return nil
}

func (prc *CollectionCore) GetCollectionCityAll(ctx context.Context, m meta.Metadata) ([]entity.City, *meta.Metadata, error) {
	var query = prc.Provider.CollectionCore.NewQuery().Col("cities").Limit(m.PerPage).Page(m.Page)
	defer query.Close()

	var result = make([]entity.City, 0)

	switch m.OrderType {
	case meta.SortAscending:
		query = query.OrderBy(m.OrderBy, collection_core.ASC)
	case meta.SortDescending:
		query = query.OrderBy(m.OrderBy, collection_core.DESC)
	}

	inf, err := query.Retrive(ctx, &result)
	if err != nil {
		return nil, nil, err
	}

	m.Page = int(inf.Meta.Page)
	m.PerPage = int(inf.Meta.PerPage)
	m.Total = int(inf.Meta.Size)
	return result, &m, nil
}

func (prc *CollectionCore) GetCollectionCityByDocumentID(ctx context.Context, documentID string) (*entity.City, error) {
	var query = prc.Provider.CollectionCore.NewQuery().Col("cities").Doc(documentID)
	defer query.Close()

	var result entity.City

	inf, err := query.Retrive(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("document %v: with reason: %v", inf.Status.String(), inf.Message)
	}

	return &result, nil
}
