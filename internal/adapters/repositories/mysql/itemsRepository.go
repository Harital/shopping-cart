package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Harital/shopping-cart/internal/core/model"
	"github.com/huandu/go-sqlbuilder"
)

const (
	cartItemTable = "cartItem"
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
		Select("id", "name", "quantity", "reservationId").
		From(cartItemTable)

	query, args := sb.Build()
	rows, selectErr := cir.db.QueryContext(ctx, query, args...)
	if selectErr != nil {
		return []model.CartItem{}, selectErr
	}

	defer rows.Close() // do not forget to close the rows. Only needed if no error

	var items []model.CartItem
	for rows.Next() {
		var singleItem model.CartItem

		// ReservationID can be null, hence the need of the sql.NullString
		var reservationId sql.NullString
		if scanErr := rows.Scan(&singleItem.Id, &singleItem.Name, &singleItem.Quantity, &reservationId); scanErr != nil {
			return []model.CartItem{}, fmt.Errorf("Scanning cart items properties --> %w", scanErr)
		}

		if reservationId.Valid {
			singleItem.ReservationId = reservationId.String
		}

		items = append(items, singleItem)
	}

	return items, nil
}

func (cir *CartItemsRepository) Add(ctx context.Context, item model.CartItem) error {

	sb := sqlbuilder.MySQL.NewInsertBuilder()
	sb.
		InsertInto(cartItemTable).
		Cols("id", "name", "quantity").
		Values(item.Id, item.Name, item.Quantity)

	// I know it´s a little nasty, but sqlbuilder does not support on dupplicated key
	// https://github.com/huandu/go-sqlbuilder/issues/15
	builder := sqlbuilder.Build("$? ON DUPLICATE KEY UPDATE quantity = quantity + $?", sb, item.Quantity)
	query, args := builder.Build()
	_, insertErr := cir.db.ExecContext(ctx, query, args...)

	if insertErr != nil {
		return fmt.Errorf("inserting items to cart --> %w", insertErr)
	}

	return nil
}

func (cir *CartItemsRepository) SetReservationId(ctx context.Context, item model.CartItem, reservationId string) error {

	sb := sqlbuilder.MySQL.NewUpdateBuilder()
	sb.Update(cartItemTable).
		Set(sb.Assign("reservationID", reservationId)).
		Where(sb.Equal("id", item.Id))

	query, args := sb.Build()
	result, updateErr := cir.db.ExecContext(ctx, query, args...)

	if updateErr != nil {
		return fmt.Errorf("cannot update reservationID --> %w", updateErr)
	}
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("cannot check rows affected when updating reservationID --> %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item id %s not found", item.Id)
	}
	return nil
}
