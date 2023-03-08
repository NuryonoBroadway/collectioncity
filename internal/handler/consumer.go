// Package handler
package handler

import (
	"context"

	"collection-squad/collection/collection-city/internal/appctx"
	"collection-squad/collection/collection-city/internal/consts"
	uContract "collection-squad/collection/collection-city/internal/ucase/contract"
	"collection-squad/collection/collection-city/pkg/awssqs"
)

// SQSConsumerHandler sqs consumer message processor handler
func SQSConsumerHandler(msgHandler uContract.MessageProcessor) awssqs.MessageProcessorFunc {
	return func(decoder *awssqs.MessageDecoder) error {
		return msgHandler.Serve(context.Background(), &appctx.ConsumerData{
			Body:        []byte(*decoder.Body),
			Key:         []byte(*decoder.MessageId),
			ServiceType: consts.ServiceTypeConsumer,
		})
	}
}
