package collection_core

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type _Payload struct {
	_client *Client `json:"-"`

	Root_Collection string        `json:"root_collection"`
	Root_Document   string        `json:"root_document"`
	Data            _Data         `json:"data"`
	Path            []_Path_Build `json:"path"`
	Limits          int32         `json:"limit"`
	Pages           int32         `json:"page"`
	Query           _Queries      `json:"query"`
	IsPagination    bool          `json:"is_pagination"`
}

// Document Ref
type DocumentRef interface {
	Col(id string) CollectionRef

	ResponseRef

	Close() error
}

func (p *_Payload) Col(id string) CollectionRef {
	p.Path = append(p.Path, _Path_Build{
		CollectionID: id,
	})
	return p
}

// Collection Ref

type CollectionRef interface {
	Doc(id string) DocumentRef
	NewDoc() DocumentRef

	QueryRef
	ResponseRef

	Close() error
}

func (p *_Payload) Doc(id string) DocumentRef {
	p.Path = append(p.Path, _Path_Build{
		DocumentID: id,
	})
	return p
}

func (p *_Payload) NewDoc() DocumentRef {
	p.Path = append(p.Path, _Path_Build{
		NewDocument: true,
	})
	return p
}

// Query Ref

type QueryRef interface {
	Where(by string, op Operator, val interface{}) QueryRef
	DateRange(field string, start time.Time, end time.Time) QueryRef
	Page(page int) QueryRef
	Limit(limit int) QueryRef
	OrderBy(by string, dir Direction) QueryRef

	Retrive(ctx context.Context, data interface{}, opts ...RetriveOption) (*ResponseInformation, error)

	Close() error
}

func (p *_Payload) DateRange(field string, start time.Time, end time.Time) QueryRef {
	p.Query.Date_Range = _Date_Range_Query{
		Field: field,
		Start: start,
		End:   end,
	}

	return p
}

func (p *_Payload) Where(by string, op Operator, val interface{}) QueryRef {
	p.Query.Filter = append(p.Query.Filter, _Filter_Query{
		By:  by,
		Op:  op.ToString(),
		Val: val,
	})

	return p
}

func (p *_Payload) Limit(limit int) QueryRef {
	if p.Pages <= 0 {
		p.Pages = 1
	}

	if limit <= 0 {
		p.Limits = 20
	}

	p.IsPagination = true
	p.Limits = int32(limit)
	return p
}

func (p *_Payload) OrderBy(by string, dir Direction) QueryRef {
	for i := 0; i < len(p.Query.Sort); i++ {
		if p.Query.Sort[i].OrderBy == by {
			return p
		}
	}

	p.Query.Sort = append(p.Query.Sort, _Sort_Query{
		OrderBy:   by,
		OrderType: dir.ToInt32(),
	})

	return p
}

func (p *_Payload) Page(page int) QueryRef {
	if page <= 0 {
		p.Pages = 1
	}

	if p.Limits <= 0 {
		p.Limits = 20
	}

	p.IsPagination = true
	p.Pages = int32(page)
	return p
}

// Response Ref

type ResponseRef interface {
	// Support Multiple Insert
	//
	// Please add method on your struct :
	//
	// Only Support Data Type:
	//	- struct
	//	- map
	//	- []struct
	//
	// Example:
	// type SomeStruct struct {
	//	ID string
	// }
	//
	// func (x *SomeStruct) RefID() string {
	// 	return x.ID
	// }
	//
	Set(ctx context.Context, data interface{}, opts ...SaveOption) (*ResponseInformation, error)
	// Only Support Data Type:
	//	- *struct
	//	- *map
	//	- *[]struct
	//	- *[]map
	//
	Retrive(ctx context.Context, data interface{}, opts ...RetriveOption) (*ResponseInformation, error)

	Close() error
}

func (p *_Payload) Retrive(ctx context.Context, data interface{}, opts ...RetriveOption) (*ResponseInformation, error) {
	req_id := uuid.New()
	if data == nil {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.InvalidArgument,
			Message: "data must be pointer",
		}, Error("collection_core_client@retrive: data must be pointer")
	}

	rType := reflect.TypeOf(data)
	if rType.Kind() != reflect.Ptr {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.InvalidArgument,
			Message: "response data must be pointer",
		}, Error("collection_core_client@retrive: response data must be pointer")
	}

	rType = rType.Elem()

	switch rType.Kind() {
	case
		reflect.Struct,
		reflect.Map,
		reflect.Slice:

		var (
			buff       = new(bytes.Buffer)
			info, _err = p._client._Retrive(ctx, p, buff)
		)
		if _err != nil {
			return info, _err
		}

		defer buff.Reset()

		if _err := json.NewDecoder(buff).Decode(&data); _err != nil {
			return &ResponseInformation{
				ReqID:   req_id,
				Status:  codes.Internal,
				Message: "failed to decode response data",
			}, Error("collection_core_client@retrive: failed to decode response data")
		}

		return info, nil
	default:
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.Internal,
			Message: "not supported type response data",
		}, Error("collection_core_client@retrive: not supported type response data")
	}
}

func (p *_Payload) Set(ctx context.Context, data interface{}, opts ...SaveOption) (*ResponseInformation, error) {
	req_id := uuid.New()
	_opts := _SaveOptions{}
	for i := 0; i < len(opts); i++ {
		if _err := opts[i](&_opts); _err != nil {
			return &ResponseInformation{
				ReqID:   req_id,
				Status:  codes.Internal,
				Message: _err.Error(),
			}, _err
		}
	}

	var (
		object = make([]_Object, 0)
		rValue = reflect.ValueOf(data)
	)

	for {
		if kind := rValue.Kind(); kind == reflect.Ptr {
			rValue = rValue.Elem()
			continue
		}
		goto _FirstStopLoop
	}
_FirstStopLoop:

	_res := ResponseInformation{
		ReqID:   req_id,
		Status:  codes.Internal,
		Message: "not supported data type",
	}

	switch rValue.Kind() {
	case reflect.Struct, reflect.Map:
		var (
			_map      = make(map[string]interface{})
			_bytes, _ = json.Marshal(data)
			_         = json.Unmarshal(_bytes, &_map)
		)

		object = append(object, _Object{
			Object: _map,
		})

	case reflect.Slice:
		args := make([]reflect.Value, 0)
		for i := 0; i < rValue.Len(); i++ {
			idx_rval := rValue.Index(i)

			for {
				if idx_kind := idx_rval.Kind(); idx_kind == reflect.Ptr {
					idx_rval = idx_rval.Elem()
					continue
				}
				goto _ElemStopLoop
			}
		_ElemStopLoop:

			idx_method := idx_rval.MethodByName("RefID")
			if idx_method.Kind() != reflect.Func {
				return &_res, ErrInvalidDataType
			}

			returns := idx_method.Call(args)
			if len(returns) <= 0 {
				return &_res, ErrInvalidDataType
			}

			ref_id, ok := returns[0].Interface().(string)
			if !ok {
				return &_res, ErrInvalidDataType
			}

			var (
				_map      = make(map[string]interface{})
				_bytes, _ = json.Marshal(idx_rval.Interface())
				_         = json.Unmarshal(_bytes, &_map)
			)

			object = append(object, _Object{
				RefID:  ref_id,
				Object: _map,
			})
		}

	default:
		return &_res, ErrInvalidDataType
	}

	p.Data.Objects = object

	if _opts.is_merge_all {
		p.Data.IsMergeAll = true
	}

	if !_opts.use_grpc && !_opts.use_pubsub {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.Internal,
			Message: "must use method between grpc and pubsub",
		}, Error("collection_core_client@save: must use method between grpc and pubsub")
	}

	if _opts.use_grpc {
		info, _err := p._client._SaveGRPC(ctx, p)
		if _err != nil {
			return info, _err
		}

		return info, nil
	}

	if p._client.pubsub_conn == nil {
		return &ResponseInformation{
			ReqID:   req_id,
			Status:  codes.Internal,
			Message: "unimplemented method pubsub",
		}, Error("collection_core_client@save: must setup configuration method pubsub")
	}

	info, _err := p._client._SavePubSub(ctx, p)
	if _err != nil {
		return info, _err
	}

	return info, nil
}

func (p *_Payload) Close() error {
	*p = _Payload{}
	return nil
}
