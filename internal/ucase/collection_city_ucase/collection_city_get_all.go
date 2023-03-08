package collection_city_ucase

import (
	"fmt"
	"net/http"

	"gitlab.privy.id/collection/collection-city/internal/appctx"
	"gitlab.privy.id/collection/collection-city/internal/provider/collection_core_provider"
	"gitlab.privy.id/collection/collection-city/internal/ucase/contract"
	"gitlab.privy.id/collection/collection-city/pkg/logger"
	"gitlab.privy.id/collection/collection-city/pkg/tracer"
	"gitlab.privy.id/collection/collection-package-core/request/meta"
	"gitlab.privy.id/collection/collection-package-core/response/errbank"
)

type GetCollectionBook struct {
	prv collection_core_provider.CollectionCoreProvider
}

func NewGetCollectionBook(prv collection_core_provider.CollectionCoreProvider) contract.UseCase {
	return &GetCollectionBook{prv: prv}
}

func (s *GetCollectionBook) EventName() string {
	return "ucase.collection_book.create"
}

func (s *GetCollectionBook) Serve(data *appctx.Data) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName(s.EventName()),
		)
		ctx            = data.Request.Context()
		m              = meta.MetadataFromURL(data.Request.URL.Query())
		response       = appctx.NewResponse()
		MessageSuccess = "Get Collection Book Success"
		MessageError   = "Get Collection Book Error"
		vErrors        = make(errbank.ValidationError, 0)
	)

	ctx = tracer.SpanStart(ctx, s.EventName())
	defer tracer.SpanFinish(ctx)

	resp, metadata, err := s.prv.GetCollectionCityAll(ctx, m)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("err get response from provider %v", err), lf...)

		vErrors = append(vErrors, errbank.FieldError{
			Field: "collection_core_get_all",
			Error: err.Error(),
		})

		return *response.WithCode(http.StatusUnprocessableEntity).WithMessage(MessageError).WithError(vErrors)
	}

	logger.InfoWithContext(ctx, fmt.Sprintf("success send payload %v", err), lf...)
	return *response.WithCode(http.StatusOK).WithMeta(metadata).WithMessage(MessageSuccess).WithData(resp)
}
