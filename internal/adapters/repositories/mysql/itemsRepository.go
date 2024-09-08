package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Harital/shopping-cart/internal/core/model"
	"github.com/huandu/go-sqlbuilder"
)

type CartItemsRepository struct {
	db *sql.DB
}

func NewCartItemsRepository(db *sql.DB) *CartItemsRepository {
	return &CartItemsRepository{db: db}
}

func (cir CartItemsRepository) Get(ctx context.Context) ([]model.CartItem, error) {

	// Simple query. An stored procedure could be used to speed up the operation.
	sb := sqlbuilder.MySQL.NewSelectBuilder()
	sb.
		Select("name", "quantity", "reservationId").
		From("cartItem")

	query, args := sb.Build()
	rows, selectErr := cir.db.QueryContext(ctx, query, args...)
	if selectErr != nil {
		return []model.CartItem{}, selectErr
	}

	defer rows.Close() // do not forget to close the rows. Only needed if no error

	var items []model.CartItem
	for rows.Next() {
		var singleItem model.CartItem

		if scanErr := rows.Scan(&singleItem.Name, &singleItem.Quantity, &singleItem.ReservationId); scanErr != nil {
			return []model.CartItem{}, fmt.Errorf("Scanning cart items properties --> %w", scanErr)
		}

		items = append(items, singleItem)
	}

	return items, nil
}

func (cir *CartItemsRepository) Add(ctx context.Context, item model.CartItem) error {

	addErr := cir.Add(ctx, item)
	if addErr != nil {
		return addErr
	}

	
	go cir.reserveItem(item)
	return nil
}

func (cir *CartItemsRepository) reserveItem(item model.CartItem) {
	//ctx, cancel := context.WithTimeout()

}
