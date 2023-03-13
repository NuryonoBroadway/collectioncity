package collection_core

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type ResponseInformation struct {
	ReqID   uuid.UUID  `json:"req_id"`
	Status  codes.Code `json:"status"`
	Meta    Meta       `json:"meta"`
	Message string     `json:"message"`
}

type Meta struct {
	Page    int32 `json:"page"`
	PerPage int32 `json:"per_page"`
	Size    int32 `json:"size"`
}

type _Object struct {
	RefID  string                 `json:"ref_id"`
	Object map[string]interface{} `json:"object"`
}

// Private

type _Data struct {
	IsMergeAll bool      `json:"is_merge_all"`
	Objects    []_Object `json:"objects"`
}

type _Path_Build struct {
	CollectionID string `json:"collection_id,omitempty"`
	DocumentID   string `json:"document_id,omitempty"`
	NewDocument  bool   `json:"new_document,omitempty"`
}

type _Queries struct {
	Sort       []_Sort_Query     `json:"sort"`
	Filter     []_Filter_Query   `json:"filter"`
	Date_Range _Date_Range_Query `json:"date_range"`
}

type _Sort_Query struct {
	OrderBy   string `json:"order_by"`
	OrderType int32  `json:"order_type"`
}

type _Filter_Query struct {
	By  string      `json:"by"`
	Op  string      `json:"op"`
	Val interface{} `json:"val"`
}

type _Date_Range_Query struct {
	Field string    `json:"field"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
