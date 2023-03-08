package collection_core_provider

import (
	"context"
	"fmt"

	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/entity"
	"gitlab.privy.id/privypass/privypass-package-core/collection_core"
	"gitlab.privy.id/privypass/privypass-package-core/request/meta"
)

func (prc *CollectionCore) CreateCollectionCityGrpc(ctx context.Context, city entity.City) error {
	var new_query = prc.Provider.CollectionCore.NewQuery()
	var query = new_query.Col("cities").Doc(city.ID)

	defer func() {
		new_query.Close()
		query.Close()
	}()

	_, err := query.Set(ctx, city, collection_core.IsMergeAll, collection_core.IsUseGRPC)
	if err != nil {
		return err
	}

	return nil
}

func (prc *CollectionCore) CreateCollectionCityPubsub(ctx context.Context, city entity.City) error {
	var new_query = prc.Provider.CollectionCore.NewQuery()
	var query = new_query.Col("cities").Doc(city.ID)

	defer func() {
		new_query.Close()
		query.Close()
	}()

	_, err := query.Set(ctx, city, collection_core.IsMergeAll, collection_core.IsUsePubSub)
	if err != nil {
		return err
	}

	return nil
}

func (prc *CollectionCore) GetCollectionCityAll(ctx context.Context, m meta.Metadata) ([]entity.City, *meta.Metadata, error) {
	var new_query = prc.Provider.CollectionCore.NewQuery()
	var query = new_query.Col("cities").Limit(m.PerPage).Page(m.Page)
	var result = make([]entity.City, 0)

	defer func() {
		new_query.Close()
		query.Close()
	}()

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
	var new_query = prc.Provider.CollectionCore.NewQuery()
	var query = new_query.Col("cities").Doc(documentID)
	var result entity.City

	defer func() {
		new_query.Close()
		query.Close()
	}()

	inf, err := query.Retrive(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("document %v: with reason: %v", inf.Status.String(), inf.Message)
	}

	return &result, nil
}
