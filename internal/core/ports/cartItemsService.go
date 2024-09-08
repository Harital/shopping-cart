package ports

import (
	"context"

	"github.com/Harital/shopping-cart/internal/core/model"
)

//go:generate mockgen -destination=../mocks/CartItemsService_mock.go -package=mocks . CartItemsService
type CartItemsService interface {
	Get(ctx context.Context) ([]model.CartItem, error)
	Add(ctx context.Context, items []model.CartItem) error
}
