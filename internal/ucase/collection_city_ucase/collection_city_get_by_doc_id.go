package collection_city_ucase

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.privy.id/collection/collection-city/internal/appctx"
	"gitlab.privy.id/collection/collection-city/internal/provider/collection_core_provider"
	"gitlab.privy.id/collection/collection-city/internal/ucase/contract"
	"gitlab.privy.id/collection/collection-city/pkg/logger"
	"gitlab.privy.id/collection/collection-city/pkg/tracer"
	"gitlab.privy.id/collection/collection-package-core/response/errbank"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetDetailCollectionBook struct {
	prv collection_core_provider.CollectionCoreProvider
}

func NewGetDetailCollectionBook(prv collection_core_provider.CollectionCoreProvider) contract.UseCase {
	return &GetDetailCollectionBook{prv: prv}
}

func (s *GetDetailCollectionBook) EventName() string {
	return "ucase.collection_book.create"
}

func (s *GetDetailCollectionBook) Serve(data *appctx.Data) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName(s.EventName()),
		)
		ctx            = data.Request.Context()
		response       = appctx.NewResponse()
		MessageSuccess = "Get Detail Collection Book Success"
		MessageError   = "Get Detail Collection Book Error"
		vErrors        = make(errbank.ValidationError, 0)
		document_id    = mux.Vars(data.Request)["document_id"]
	)

	ctx = tracer.SpanStart(ctx, s.EventName())
	defer tracer.SpanFinish(ctx)

	if len(document_id) == 0 {
		err := errors.New("param is empty")
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("err get response from provider %v", err), lf...)

		vErrors = append(vErrors, errbank.FieldError{
			Field: "document_id",
			Error: err.Error(),
		})

		return *response.WithCode(http.StatusNotFound).WithMessage(MessageError).WithError(vErrors)
	}

	resp, err := s.prv.GetCollectionCityByDocumentID(ctx, document_id)
	if err != nil {
		stat := status.Convert(err)
		switch stat.Code() {
		case codes.Internal:
			tracer.SpanError(ctx, err)
			logger.ErrorWithContext(ctx, fmt.Sprintf("err get response from provider %v", err), lf...)

			return *response.WithCode(http.StatusInternalServerError).WithError(stat.Message())

		default:
			tracer.SpanError(ctx, err)
			logger.ErrorWithContext(ctx, fmt.Sprintf("err get response from provider %v", err), lf...)

			vErrors = append(vErrors, errbank.FieldError{
				Field: "collection_core_get_all",
				Error: stat.Message(),
			})

			return *response.WithCode(http.StatusUnprocessableEntity).WithMessage(MessageError).WithError(vErrors)
		}

	}

	logger.InfoWithContext(ctx, fmt.Sprintf("success send payload %v", err), lf...)
	return *response.WithCode(http.StatusOK).WithMessage(MessageSuccess).WithData(resp)
}
