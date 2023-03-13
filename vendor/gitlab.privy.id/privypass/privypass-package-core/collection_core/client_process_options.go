package collection_core

const (
	ErrInvalidDataType         = Error("collection_core_client@save: invalid data type")
	ErrUseMultipleMethodOption = Error("collection_core_client@save: cannot use both, choose one between grpc and pubsub")
)

type _SaveOptions struct {
	is_merge_all bool
	use_pubsub   bool
	use_grpc     bool
}

type SaveOption func(op *_SaveOptions) error

var (
	IsMergeAll SaveOption = func(op *_SaveOptions) error {
		op.is_merge_all = true
		return nil
	}

	IsUsePubSub SaveOption = func(op *_SaveOptions) error {
		if op.use_grpc {
			return ErrUseMultipleMethodOption
		}
		op.use_pubsub = true
		return nil
	}

	IsUseGRPC SaveOption = func(op *_SaveOptions) error {
		if op.use_pubsub {
			return ErrUseMultipleMethodOption
		}
		op.use_grpc = true
		return nil
	}
)

type _RetriveOptions struct{}

type RetriveOption func(op *_RetriveOptions) error
