package collection_core

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"gitlab.privy.id/privypass/privypass-package-core/collection_core/pb"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	ctx    context.Context
	config *_Config

	grpc_conn   *grpc.ClientConn
	pubsub_conn *pubsub.Client

	publisher              *gPublisher
	collection_grpc_client pb.ServiceCollectionCoreClient
}

func NewClient(ctx context.Context, opts ...ConfigOption) (*Client, error) {
	c := &Client{
		ctx:    ctx,
		config: &_Config{},
	}

	for i := 0; i < len(opts); i++ {
		if _err := opts[i](c.config); _err != nil {
			return nil, _err
		}
	}

	var (
		IsHasCredentialFile = len(c.config.pubsub.credential_file) > 0
		IsHasProjectID      = len(c.config.pubsub.project_id) > 0
		IsHasProjectName    = len(c.config.pubsub.project_name) > 0
		IsHasTopic          = len(c.config.pubsub.topic) > 0

		IsUsePubSub = (IsHasCredentialFile && IsHasProjectID && IsHasProjectName && IsHasTopic)
	)

	if IsUsePubSub {
		// Open Connection Google PUB/SUB
		pubsub_client, err := pubsub.NewClient(
			ctx,
			c.config.pubsub.project_id,
			option.WithCredentialsFile(c.config.pubsub.credential_file),
		)
		if err != nil {
			c.Close()
			return nil, err
		}

		var (
			topic        = pubsub_client.Topic(c.config.pubsub.topic)
			exists, _err = topic.Exists(ctx)
		)
		if _err != nil {
			c.Close()
			return nil, err
		}

		if !exists {
			c.Close()
			return nil, fmt.Errorf("topic '%s' not registered", c.config.pubsub.topic)
		}

		c.pubsub_conn = pubsub_client
		c.publisher = NewGPublisher(c.pubsub_conn)
	}

	// Open Client Long Connection GRPC Server

	grpc_client, err := grpc.DialContext(
		ctx,
		c.config.grpc.address,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		c.Close()
		return nil, err
	}

	c.grpc_conn = grpc_client
	c.collection_grpc_client = pb.NewServiceCollectionCoreClient(grpc_client)

	return c, nil
}

func (c *Client) Close() {
	if c.pubsub_conn != nil {
		c.pubsub_conn.Close()
	}

	if c.grpc_conn != nil {
		c.grpc_conn.Close()
	}

	if c.collection_grpc_client != nil {
		c.collection_grpc_client = nil
	}

	if c.publisher != nil {
		c.publisher = nil
	}
}

func (c *Client) NewQuery() DocumentRef {
	var (
		p = _Payload{
			_client: c,

			Root_Collection: c.config.firestore.root_collection_id,
			Root_Document:   c.config.firestore.root_document_id,
			Data: _Data{
				Objects: make([]_Object, 0),
			},
			Path: make([]_Path_Build, 0),
			Query: _Queries{
				Sort:   make([]_Sort_Query, 0),
				Filter: make([]_Filter_Query, 0),
			},
		}
	)

	return &p
}
