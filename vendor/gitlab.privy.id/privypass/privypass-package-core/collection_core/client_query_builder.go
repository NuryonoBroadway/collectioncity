package collection_core

import (
	"encoding/json"
	"fmt"

	"gitlab.privy.id/privypass/privypass-package-core/collection_core/pb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Client) _Builder(p *_Payload, is_retrive bool) (*pb.PayloadProto, error) {
	pp := pb.PayloadProto{
		RootCollection: p.Root_Collection,
		RootDocument:   p.Root_Document,
		Data: &pb.DataProto{
			IsMergeAll: p.Data.IsMergeAll,
			Objects:    make([]*pb.ObjectProto, len(p.Data.Objects)),
		},
		Path: make([]*pb.PathProto, len(p.Path)),
		Query: &pb.QueryProto{
			Sort:   make([]*pb.SortProto, len(p.Query.Sort)),
			Filter: make([]*pb.FilterProto, len(p.Query.Filter)),
			DateRange: &pb.DateRangeProto{
				Field: p.Query.Date_Range.Field,
				Start: timestamppb.New(p.Query.Date_Range.Start),
				End:   timestamppb.New(p.Query.Date_Range.End),
			},
		},
		Page:         p.Pages,
		Limit:        p.Limits,
		IsPagination: p.IsPagination,
	}

	for i := 0; i < len(p.Data.Objects); i++ {
		_obj, _err := structpb.NewStruct(p.Data.Objects[i].Object)
		if _err != nil {
			_method := "save"
			if is_retrive {
				_method = "retrive"
			}
			return nil, Error(fmt.Sprintf("collection_core_client@%s: failed to mapping object data", _method))
		}

		pp.Data.Objects[i] = &pb.ObjectProto{
			RefId:  p.Data.Objects[i].RefID,
			Object: _obj,
		}
	}

	for i := 0; i < len(p.Path); i++ {
		pp.Path[i] = &pb.PathProto{
			CollectionId: p.Path[i].CollectionID,
			DocumentId:   p.Path[i].DocumentID,
			NewDocument:  p.Path[i].NewDocument,
		}
	}

	for i := 0; i < len(p.Query.Sort); i++ {
		var _sort_type pb.SortTypeProto
		switch p.Query.Sort[i].OrderType {
		case 1:
			_sort_type = pb.SortTypeProto_Asc
		case 2:
			_sort_type = pb.SortTypeProto_Desc
		default:
			_sort_type = pb.SortTypeProto_None
		}

		pp.Query.Sort[i] = &pb.SortProto{
			OrderBy:   p.Query.Sort[i].OrderBy,
			OrderType: _sort_type,
		}
	}

	for i := 0; i < len(p.Query.Filter); i++ {
		_val_bytes, _err := json.Marshal(p.Query.Filter[i].Val)
		if _err != nil {
			_method := "save"
			if is_retrive {
				_method = "retrive"
			}
			return nil, Error(fmt.Sprintf("collection_core_client@%s: failed to mapping filter data", _method))
		}

		pp.Query.Filter[i] = &pb.FilterProto{
			By:  p.Query.Filter[i].By,
			Op:  p.Query.Filter[i].Op,
			Val: _val_bytes,
		}
	}

	return &pp, nil
}
