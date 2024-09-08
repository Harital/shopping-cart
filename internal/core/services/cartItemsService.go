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

func (cis CartItemsService) Get(ctx context.Context) ([]model.CartItem, error) {
	return cis.repo.Get(ctx)
}

func (cis0 *CartItemsService) Add(ctx context.Context, items []model.CartItem) error {
	return nil
}
