package mysql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Harital/shopping-cart/internal/core/model"
	"github.com/stretchr/testify/assert"
)

const (
	cartItemsGetQuery = "SELECT name, quantity, reservationId FROM cartItem"
)

// Always useful to add mocks into the struct, so it´s easier to add more mocks in the future
// No additional parameters are needed.
type mocks struct {
	sql sqlmock.Sqlmock
}

// We use Gherkin notation for the tests
func Test_GetCartItems_GivenInitializedRepository(t *testing.T) {
	randomError := errors.New("random error")

	type want struct {
		err   error
		items []model.CartItem
	}

	tests := []struct {
		name  string
		mocks func(m mocks)
		want  want
	}{
		{
			name: "WhenGetAndErrorInSelect_ThenError",
			mocks: func(m mocks) {
				m.sql.
					ExpectQuery(cartItemsGetQuery).
					WillReturnError(randomError)
			},
			want: want{
				err:   randomError,
				items: []model.CartItem{},
			},
		}, {
			name: "WhenGetAndErrorInScan_ThenError",
			mocks: func(m mocks) {
				// We need to call NewRows in every test 
				m.sql.
					ExpectQuery(cartItemsGetQuery).
					WillReturnRows(sqlmock.NewRows([]string{
						"name", "quantity", "reservationId",
					}).
						AddRow("pants", nil, "reservationId1"))
			},
			want: want{
				err:   randomError,
				items: []model.CartItem{},
			},
		}, {
			name: "WhenGetAnd1ItemReturn_ThenOK",
			mocks: func(m mocks) {
				m.sql.
					ExpectQuery(cartItemsGetQuery).
					WillReturnRows(sqlmock.NewRows([]string{
						"name", "quantity", "reservationId",
					}).
						AddRow("pants", 1, "reservationId1"))
			},
			want: want{
				err:   nil,
				items: []model.CartItem{
					{Name: "pants", Quantity: 1, ReservationId: "reservationId1"},
				},
			},
		}, {
			name: "WhenGetAndSeveralItemsReturn_ThenOK",
			mocks: func(m mocks) {
				m.sql.
					ExpectQuery(cartItemsGetQuery).
					WillReturnRows(sqlmock.NewRows([]string{
						"name", "quantity", "reservationId",
					}).
						AddRow("bottle", 10, "reservationId5",).
						AddRow("shirt", 2, "reservationId2",))
			},
			want: want{
				err:   nil,
				items: []model.CartItem{
					{Name: "bottle", Quantity: 10, ReservationId: "reservationId5"},
					{Name: "shirt", Quantity: 2, ReservationId: "reservationId2"},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("Error when creating the mock: %v", err)
			}
			m := mocks{sql: dbMock}
			defer db.Close()

			tc.mocks(m)

			r := NewCartItemsRepository(db)

			items, getErr := r.Get(context.TODO())
			if tc.want.err != nil {
				// helps being agnostic with the error message, as it can change and wrongfully break the tests
				// If an error comprobation is needed, assert.ErrorIs or ErrorAs can be used.
				assert.Error(t, getErr)
			} else {
				assert.NoError(t, getErr)
			}
			assert.Equal(t, tc.want.items, items)
		})

	}
}
