package services

import (
	"context"
	"net/http"
	"time"

	"github.com/Harital/shopping-cart/internal/core/model"
	"github.com/Harital/shopping-cart/internal/core/ports"
)

type CartItemsService struct {
	repo ports.CartItemsRepository
	reserverTimeout time.Duration
}

func NewCartItemsService(repo ports.CartItemsRepository, reserverTimeout time.Duration) *CartItemsService {
	return &CartItemsService{
		repo: repo, 
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
	if addErr := cis.repo.Add(ctx, item); addErr != nil {
		return addErr
	}

	go cis.ReserveItem(ctx, item)

	return nil
}

func (cis *CartItemsService) ReserveItem(parentContext context.Context, item model.CartItem) {
	// Create a specific context with a big timeout for this operation
	_, cancel := context.WithTimeout(context.Background(), cis.reserverTimeout)
	defer cancel()


}