package bootstrap

import (
	"context"

	"gitlab.privy.id/privypass/privypass-oauth2-core-se/internal/appctx"
	"gitlab.privy.id/privypass/privypass-oauth2-core-se/pkg/logger"
	"gitlab.privy.id/privypass/privypass-package-core/collection_core"
)

func RegistryCollectionCore(conf *appctx.Config) *collection_core.Client {
	ctx := context.Background()
	var (
		collectionCoreInitializeNil  = `collection core cannot connect, please check your config or network`
		logFirestoreRootCollectionID = "firebase_root_collection_id"
		logFirestoreRootDocumentID   = "firebase_root_document_id"
		logGrpcAdress                = "grpc_address"
		logPubsubProjectID           = "pubsub_project_id"
		logPubsubProjectName         = "pubsub_project_name"
		logPubsubCredentialFile      = "pubsub_credential_file"
		logPubsubTopic               = "pubsub_topic"
	)

	var (
		FirestoreRootCollectionID = conf.Firestore.RootCollectionID
		FirestoreRootDocumentID   = conf.Firestore.RootDocumentID
		GRPCAddress               = conf.App.GrpcPort
		PubsubProjectID           = conf.PubSub.ProjectID
		PubsubProjectName         = conf.PubSub.ProjectName
		PubsubCredentialFile      = conf.PubSub.AccountPath
		PubsubTopic               = conf.PubSub.Topic
	)

	lf := []logger.Field{
		logger.Any(logFirestoreRootCollectionID, conf.Firestore.RootCollectionID),
		logger.Any(logFirestoreRootDocumentID, conf.Firestore.RootDocumentID),
		logger.Any(logGrpcAdress, conf.App.GrpcPort),
		logger.Any(logPubsubProjectID, conf.PubSub.ProjectID),
		logger.Any(logPubsubProjectName, conf.PubSub.ProjectName),
		logger.Any(logPubsubCredentialFile, conf.PubSub.AccountPath),
		logger.Any(logPubsubTopic, conf.PubSub.Topic),
	}

	cfg := []collection_core.ConfigOption{
		collection_core.WithConfig_Firestore(FirestoreRootCollectionID, FirestoreRootDocumentID),
		collection_core.WithConfig_GRPC(GRPCAddress),
		collection_core.WithConfig_PubSub(PubsubProjectID, PubsubProjectName, PubsubCredentialFile, PubsubTopic),
	}

	client, err := collection_core.NewClient(ctx, cfg...)
	if err != nil {
		logger.Fatal(collectionCoreInitializeNil, lf...)
	}

	return client
}
