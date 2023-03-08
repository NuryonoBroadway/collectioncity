package provider

import "gitlab.privy.id/privypass/privypass-package-core/collection_core"

type Provider struct {
	CollectionCore *collection_core.Client
}

func NewProvider(ColllectionCore *collection_core.Client) *Provider {
	return &Provider{CollectionCore: ColllectionCore}
}
