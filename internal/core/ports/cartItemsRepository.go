package ports

import (
	"context"

	"github.com/Harital/shopping-cart/internal/core/model"
)

//go:generate mockgen -destination=../mocks/CartItemsRepository_mock.go -package=mocks . CartItemsRepository
type CartItemsRepository interface {
	Get(ctx context.Context) ([]model.CartItem, error)
	Add(ctx context.Context, item model.CartItem) (error)
}
