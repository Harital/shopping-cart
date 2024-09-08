package services

import (
	"context"

	"github.com/Harital/shopping-cart/internal/core/model"
	"github.com/Harital/shopping-cart/internal/core/ports"
)

type CartItemsService struct {
	repo ports.CartItemsRepository
}

func NewCartItemsService(repo ports.CartItemsRepository) *CartItemsService {
	return &CartItemsService{repo: repo}
}

// There should be a parameter in the get of the cart that we are getting. 
// On top of that, we should check that we have the needed permissions to access this cart
// For the sake of simplicity, not implemented. Always accessing the same cart.
func (cis CartItemsService) Get(ctx context.Context) ([]model.CartItem, error) {
	return cis.repo.Get(ctx)
}

func (cis0 *CartItemsService) Add(ctx context.Context, items []model.CartItem) error {
	return nil
}
