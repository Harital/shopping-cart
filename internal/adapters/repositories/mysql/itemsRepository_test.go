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

var (
	randomError = errors.New("random error")
)

// Always useful to add CartItemRepoMocks into the struct, so it´s easier to add more CartItemRepoMocks in the future
// No additional parameters are needed.
type CartItemRepoMocks struct {
	sql sqlmock.Sqlmock
}

// We use Gherkin notation for the tests
func Test_GetCartItems_GivenInitializedRepository(t *testing.T) {

	type want struct {
		err   error
		items []model.CartItem
	}

	tests := []struct {
		name  string
		mocks func(m CartItemRepoMocks)
		want  want
	}{
		{
			name: "WhenGetAndErrorInSelect_ThenError",
			mocks: func(m CartItemRepoMocks) {
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
			mocks: func(m CartItemRepoMocks) {
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
			mocks: func(m CartItemRepoMocks) {
				m.sql.
					ExpectQuery(cartItemsGetQuery).
					WillReturnRows(sqlmock.NewRows([]string{
						"name", "quantity", "reservationId",
					}).
						AddRow("pants", 1, "reservationId1"))
			},
			want: want{
				err: nil,
				items: []model.CartItem{
					{Name: "pants", Quantity: 1, ReservationId: "reservationId1"},
				},
			},
		}, {
			name: "WhenGetAndSeveralItemsReturn_ThenOK",
			mocks: func(m CartItemRepoMocks) {
				m.sql.
					ExpectQuery(cartItemsGetQuery).
					WillReturnRows(sqlmock.NewRows([]string{
						"name", "quantity", "reservationId",
					}).
						AddRow("bottle", 10, "reservationId5").
						AddRow("shirt", 2, "reservationId2"))
			},
			want: want{
				err: nil,
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
			m := CartItemRepoMocks{sql: dbMock}
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

func Test_AddCartItems_GivenInitializedRepository(t *testing.T) {
	insertQuery := `INSERT INTO cartItem (id, name, quantity) 
		VALUES (?, ?, ?) ON DUPLICATED KEY UPDATE quantity = quantity + ?`

	randomCartItem := model.CartItem{
		Id:       "1",
		Name:     "screen",
		Quantity: 2,
	}

	type input struct {
		item model.CartItem
	}
	type want struct {
		err error
	}

	tests := []struct {
		name  string
		in    input
		mocks func(m CartItemRepoMocks)
		want  want
	}{
		{
			name: "WhenAddItemAndInsertError_ThenError",
			in: input{
				item: randomCartItem,
			},
			mocks: func(m CartItemRepoMocks) {
				m.sql.
					ExpectExec(insertQuery).
					WithArgs(randomCartItem.Id, randomCartItem.Name, randomCartItem.Quantity, randomCartItem.Quantity).
					WillReturnError(errors.New("insert error"))
			},
			want: want{
				err: randomError,
			},
		}, {
			name: "WhenAddItemAndInsertOK_ThenOK",
			in: input{
				item: randomCartItem,
			},
			mocks: func(m CartItemRepoMocks) {
				m.sql.
					ExpectExec(insertQuery).
					WithArgs(randomCartItem.Id, randomCartItem.Name, randomCartItem.Quantity, randomCartItem.Quantity).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("Error when creating the mock: %v", err)
			}
			m := CartItemRepoMocks{sql: dbMock}
			defer db.Close()

			tc.mocks(m)

			r := NewCartItemsRepository(db)

			addErr := r.Add(context.TODO(), tc.in.item)

			if tc.want.err != nil {
				assert.Error(t, addErr)
			} else {
				assert.NoError(t, addErr)
			}
		})
	}

}

func Test_AddReserveationId_GivenInitializedRepository(t *testing.T) {
	updateQuery := `UPDATE cartItem SET reservationID = ? WHERE id = ?`

	randomCartItem := model.CartItem{
		Id:       "1",
		Name:     "screen",
		Quantity: 2,
	}
	randomReservationID := "randomResvId"

	type input struct {
		item          model.CartItem
		reservationId string
	}
	type want struct {
		err error
	}

	tests := []struct {
		name  string
		in    input
		mocks func(m CartItemRepoMocks)
		want  want
	}{
		{
			name: "WhenSetReservationIdAndErrorInQuery_ThenError",
			in: input{
				item:          randomCartItem,
				reservationId: randomReservationID,
			},
			mocks: func(m CartItemRepoMocks) {
				m.sql.ExpectExec(updateQuery).
					WithArgs(randomReservationID, randomCartItem.Id).
					WillReturnError(errors.New("insert error"))
			},
			want: want{
				err: randomError,
			},
		}, {
			name: "WhenSetReservationIdAndNoRowsAffected_ThenError",
			in: input{
				item:          randomCartItem,
				reservationId: randomReservationID,
			},
			mocks: func(m CartItemRepoMocks) {
				m.sql.ExpectExec(updateQuery).
					WithArgs(randomReservationID, randomCartItem.Id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			want: want{
				err: randomError,
			},
		}, {
			name: "WhenSetReservationIdAndOK_ThenOK",
			in: input{
				item:          randomCartItem,
				reservationId: randomReservationID,
			},
			mocks: func(m CartItemRepoMocks) {
				m.sql.ExpectExec(updateQuery).
					WithArgs(randomReservationID, randomCartItem.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("Error when creating the mock: %v", err)
			}
			m := CartItemRepoMocks{sql: dbMock}
			defer db.Close()

			tc.mocks(m)

			r := NewCartItemsRepository(db)

			addErr := r.SetReservationId(context.TODO(), tc.in.item, tc.in.reservationId)

			if tc.want.err != nil {
				assert.Error(t, addErr)
			} else {
				assert.NoError(t, addErr)
			}
		})
	}
}
