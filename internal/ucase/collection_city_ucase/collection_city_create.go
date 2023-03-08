package collection_city_ucase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gitlab.privy.id/collection/collection-city/internal/appctx"
	"gitlab.privy.id/collection/collection-city/internal/consts"
	"gitlab.privy.id/collection/collection-city/internal/entity"
	"gitlab.privy.id/collection/collection-city/internal/provider/collection_core_provider"
	"gitlab.privy.id/collection/collection-city/internal/ucase/contract"
	"gitlab.privy.id/collection/collection-city/pkg/logger"
	"gitlab.privy.id/collection/collection-city/pkg/tracer"
	"gitlab.privy.id/collection/collection-package-core/response/errbank"
)

type CreateCollectionBook struct {
	prv collection_core_provider.CollectionCoreProvider
}

func NewCreateCollectionBook(prv collection_core_provider.CollectionCoreProvider) contract.UseCase {
	return &CreateCollectionBook{prv: prv}
}

func (s *CreateCollectionBook) EventName() string {
	return "ucase.collection_book.create"
}

func (s *CreateCollectionBook) Serve(data *appctx.Data) appctx.Response {
	var (
		lf = logger.NewFields(
			logger.EventName(s.EventName()),
		)
		ctx            = data.Request.Context()
		response       = appctx.NewResponse()
		MessageSuccess = "Get Collection Book Success"
		MessageError   = "Get Collection Book Error"
		vErrors        = make(errbank.ValidationError, 0)
	)

	ctx = tracer.SpanStart(ctx, s.EventName())
	defer tracer.SpanFinish(ctx)

	var payload entity.CityRequest
	err := json.NewDecoder(data.Request.Body).Decode(&payload)
	if err != nil {
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("err parsing payload %v", err), lf...)

		return *response.WithCode(http.StatusUnprocessableEntity).WithMessage(err)
	}

	if err := payload.Validate(ctx); err != nil {
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("err validate payload %v", err), lf...)
		errors := consts.ToValidationError(err)

		return *response.WithCode(http.StatusUnprocessableEntity).WithMessage(MessageError).WithError(errors)
	}

	var req = entity.City{
		ID:         uuid.NewString(),
		Name:       payload.Name,
		State:      payload.State,
		Country:    payload.Country,
		Capital:    payload.Capital,
		Population: payload.Population,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	switch {
	case payload.IsGRPC:
		err := s.prv.CreateCollectionCityGrpc(ctx, req)
		if err != nil {
			tracer.SpanError(ctx, err)
			logger.ErrorWithContext(ctx, fmt.Sprintf("err send payload %v", err), lf...)
			vErrors = append(vErrors, errbank.FieldError{
				Field: "collection_core_create_grpc",
				Error: err.Error(),
			})

			return *response.WithCode(http.StatusUnprocessableEntity).WithMessage(MessageError).WithError(vErrors)
		}
	case payload.IsPubsub:
		err := s.prv.CreateCollectionCityPubsub(ctx, req)
		if err != nil {
			tracer.SpanError(ctx, err)
			logger.ErrorWithContext(ctx, fmt.Sprintf("err send payload %v", err), lf...)
			vErrors = append(vErrors, errbank.FieldError{
				Field: "collection_core_create_grpc",
				Error: err.Error(),
			})

			return *response.WithCode(http.StatusUnprocessableEntity).WithMessage(MessageError).WithError(vErrors)
		}
	default:
		tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("err send payload %v", err), lf...)
		vErrors = append(vErrors, errbank.FieldError{
			Field: "collection_core_create",
			Error: errors.New("create type is invalid").Error(),
		})

		return *response.WithCode(http.StatusUnprocessableEntity).WithMessage(MessageError).WithError(vErrors)
	}

	logger.InfoWithContext(ctx, fmt.Sprintf("success send payload %v", err), lf...)
	return *response.WithCode(http.StatusCreated).WithMessage(MessageSuccess).WithData(req)
}
