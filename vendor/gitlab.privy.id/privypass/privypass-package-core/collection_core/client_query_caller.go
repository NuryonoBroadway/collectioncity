package collection_core

import (
	"context"
	"encoding/json"
	"io"

	"github.com/google/uuid"
	"gitlab.privy.id/privypass/privypass-package-core/collection_core/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) _Retrive(ctx context.Context, p *_Payload, _data io.Writer) (*ResponseInformation, error) {
	req_id := uuid.New()
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), DefaultRequestTimeout)
		defer func() {
			cancel()
		}()
	}

	if p == nil {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.InvalidArgument,
			Message: "payload is empty",
		}, Error("collection_core_client@retrive: payload is empty")
	}

	if _data == nil {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.InvalidArgument,
			Message: "preparation to decode data is failed",
		}, Error("collection_core_client@retrive: preparation to decode data is failed")
	}

	pp, _err := c._Builder(p, true)
	if _err != nil {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.InvalidArgument,
			Message: "build payload data is failed",
		}, _err
	}

	result, _err := c.collection_grpc_client.Retrive(ctx, &pb.RetriveRequestProto{
		ReqId:   req_id.String(),
		Payload: pp,
	})
	if _err != nil {
		_status := status.Convert(_err)
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  _status.Code(),
			Message: _status.Message(),
		}, _err
	}

	if n, _ := _data.Write(result.Data); n < len(result.Data) {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.Internal,
			Message: "decode data is failed",
		}, Error("collection_core_client@retrive: decode data is failed")
	}

	return &ResponseInformation{
		ReqID:  req_id,
		Status: codes.OK,
		Meta: Meta{
			Page:    result.Meta.Page,
			PerPage: result.Meta.PerPage,
			Size:    result.Meta.Size,
		},
		Message: "success retrive data",
	}, nil
}

func (c *Client) _SaveGRPC(ctx context.Context, p *_Payload) (*ResponseInformation, error) {
	req_id := uuid.New()
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), DefaultRequestTimeout)
		defer func() {
			cancel()
		}()
	}

	pp, _err := c._Builder(p, false)
	if _err != nil {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.InvalidArgument,
			Message: "build payload data is failed",
		}, _err
	}

	if _, _err := c.collection_grpc_client.Save(ctx, &pb.SaveRequestProto{
		ReqId:   req_id.String(),
		Payload: pp,
	}); _err != nil {
		_status := status.Convert(_err)
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  _status.Code(),
			Message: _status.Message(),
		}, _err
	}

	return &ResponseInformation{
		ReqID:   req_id,
		Status:  codes.OK,
		Message: "success save data via grpc",
	}, nil
}

func (c *Client) _SavePubSub(ctx context.Context, p *_Payload) (*ResponseInformation, error) {
	req_id := uuid.New()
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), DefaultRequestTimeout)
		defer func() {
			cancel()
		}()
	}

	_jbytes, _err := json.Marshal(p)
	if _err != nil {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.InvalidArgument,
			Message: "encode data is failed",
		}, Error("collection_core_client@save_pubsub: encode data is failed")
	}

	if _err := c.publisher.Publish(ctx, &Message{
		Topic: c.config.pubsub.topic,
		Data:  _jbytes,
	}); _err != nil {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.Internal,
			Message: "send data to pubsub is failed",
		}, Error("collection_core_client@save_pubsub: send data to pubsub is failed, got err: " + _err.Error())
	}

	return &ResponseInformation{
		ReqID:   req_id,
		Status:  codes.OK,
		Message: "success submited data via pubsub",
	}, nil
}
