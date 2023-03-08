package entity

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type City struct {
	ID         string    `json:"id,omitempty" firestore:"id,omitempty"`
	Name       string    `json:"name,omitempty" firestore:"name,omitempty"`
	State      string    `json:"state,omitempty" firestore:"state,omitempty"`
	Country    string    `json:"country,omitempty" firestore:"country,omitempty"`
	Capital    bool      `json:"capital" firestore:"capital"`
	Population int       `json:"population,omitempty" firestore:"population"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" firestore:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at" firestore:"deleted_at,omitempty"`
}

type CityRequest struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	Country    string `json:"country"`
	Capital    bool   `json:"capital"`
	Population int    `json:"population"`
	IsGRPC     bool   `json:"is_grpc,omitempty"`
	IsPubsub   bool   `json:"is_pubsub,omitempty"`
}

func (c CityRequest) Validate(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.State, validation.Required),
		validation.Field(&c.Country, validation.Required),
		validation.Field(&c.Population, validation.Required),
	)
}
