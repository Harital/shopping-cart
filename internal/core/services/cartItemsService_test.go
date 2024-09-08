package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Harital/shopping-cart/internal/core/mocks"
	"github.com/Harital/shopping-cart/internal/core/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Always useful to add mocks into the struct, so it´s easier to add more mocks in the future
// No additional parameters are needed.
type cartItemsServiceMocks struct {
	repo *mocks.MockCartItemsRepository
}

func Test_GetCartItemsService_GivenCartItemsServiceCreated(t *testing.T) {
	ctx := context.Background()
	randomError := errors.New("random error")
	sampleCartItems := []model.CartItem{
		{Name: "pants", Quantity: 1, ReservationId: "reservationId1"},
	}
	type want struct {
		err   error
		items []model.CartItem
	}
	tests := []struct {
		name  string
		mocks func(m cartItemsServiceMocks)
		want  want
	}{
		{
			name: "WhenGetAndError_ThenError",
			mocks: func(m cartItemsServiceMocks) {
				m.repo.EXPECT().
					Get(ctx).
					Return([]model.CartItem{}, randomError)
			},
			want: want{
				err:   randomError,
				items: []model.CartItem{},
			},
		}, {
			name: "WhenGetAndOK_ThenOK",
			mocks: func(m cartItemsServiceMocks) {
				m.repo.EXPECT().
					Get(ctx).
					Return(sampleCartItems, nil)
			},
			want: want{
				err:   nil,
				items: sampleCartItems,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			m := cartItemsServiceMocks{
				repo: mocks.NewMockCartItemsRepository(mockCtrl),
			}
			tc.mocks(m)

			svc := NewCartItemsService(m.repo)

			items, getErr := svc.Get(ctx)
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
