package services

import (
	"context"
	"time"

	"github.com/Harital/shopping-cart/internal/core/model"
	"github.com/Harital/shopping-cart/internal/core/ports"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

const (
	reserverEndpoint = "/reserve"
)

type CartItemsService struct {
	repo            ports.CartItemsRepository
	reserverHost    string
	reserverTimeout time.Duration
}

func NewCartItemsService(repo ports.CartItemsRepository, reserverHost string, reserverTimeout time.Duration) *CartItemsService {
	return &CartItemsService{
		repo:            repo,
		reserverHost:    reserverHost,
		reserverTimeout: reserverTimeout,
	}
}

// There should be a parameter in the get of the cart that we are getting.
// On top of that, we should check that we have the needed permissions to access this cart
// For the sake of simplicity, not implemented. Always accessing the same cart.
func (cis CartItemsService) Get(ctx context.Context) ([]model.CartItem, error) {
	return cis.repo.Get(ctx)
}

func (cis *CartItemsService) Add(ctx context.Context, item model.CartItem) error {
	// first we add the item to the database. Afterwards we reserve the item in background
	if addErr := cis.repo.Add(ctx, item); addErr != nil {
		return addErr
	}

	go cis.ReserveItem(ctx, item)

	return nil
}

func (cis *CartItemsService) ReserveItem(parentContext context.Context, item model.CartItem) {
	// Create a specific context with a big timeout for this operation
	reqCtx, cancel := context.WithTimeout(context.Background(), cis.reserverTimeout)
	defer cancel()

	// Authentication should be also set in this endpoint
	var reservationResponse model.ItemReservationResponse
	response, reservationErr := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(model.NewItemReservationRequest(item)).
		SetResult(&reservationResponse).
		SetContext(reqCtx).
		Post(cis.reserverHost + reserverEndpoint)

		log.Error().Msg("patata")
	if reservationErr != nil {
		log.
			Error().
			Err(reservationErr).
			Str("itemId", item.Id).
			Str("itemName", item.Name).
			Msg("While reserving item")
		return
	}

	if response.StatusCode() != 200 {
		log.
			Error().
			Str("httpResponse", response.String()).
			Msg("bad http response while reserving the item")
		return
	}

	// update the database
	setResvIdErr := cis.repo.SetReservationId(reqCtx, item, reservationResponse.ReservationId)
	if setResvIdErr != nil {
		log.
			Error().
			Err(setResvIdErr).
			Str("itemId", item.Id).
			Str("itemName", item.Name).
			Str("reservationId", reservationResponse.ReservationId).
			Msg("while storing resrvation id")
	}
}
