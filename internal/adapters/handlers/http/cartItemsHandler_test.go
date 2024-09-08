package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Harital/shopping-cart/internal/core/mocks"
	"github.com/Harital/shopping-cart/internal/core/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type CartItemHandlerMocks struct {
	svc *mocks.MockCartItemsService
}

func Test_GetCartItems_GivenInitializedHandler(t *testing.T) {
	internalError := errors.New("Internal Error")
	type want struct {
		httpCode int
		body     string
	}

	tests := []struct {
		name  string
		mocks func(m CartItemHandlerMocks)
		want  want
	}{
		{
			name: "WhenGetCartItemsAndError_ThenErrorIsReturned",
			mocks: func(m CartItemHandlerMocks) {
				m.svc.EXPECT().Get(gomock.Any()).
					Return([]model.CartItem{}, internalError)
			},
			want: want{
				httpCode: 500,
				body:     `{"version":"1.0.0","Message":"Error getting the cart items"}`,
			},
		}, {
			name: "WhenGetCartItemsAndOK_ThenItemsAreRetrieved",
			mocks: func(m CartItemHandlerMocks) {
				m.svc.EXPECT().Get(gomock.Any()).
					Return([]model.CartItem{
						{Name: "bottle", Quantity: 10, ReservationId: "reservationId5"},
						{Name: "mouse", Quantity: 4, ReservationId: "mouseReservationId"},
					}, nil)
			},
			want: want{
				httpCode: 200,
				body:     `{"Version":"1.0.0","Items":[{"name":"bottle","quantiy":10,"reservationId":"reservationId5"},{"name":"mouse","quantiy":4,"reservationId":"mouseReservationId"}]}`,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			respRecorder := httptest.NewRecorder()
			router := gin.Default()
			rg := router.Group("/shopping-cart/v1")

			m := CartItemHandlerMocks{
				svc: mocks.NewMockCartItemsService(mockCtrl),
			}

			tc.mocks(m)

			hndl := NewCartItemsHandler(rg, m.svc)
			hndl.Register()

			req := httptest.NewRequest(http.MethodGet, "/shopping-cart/v1/items", nil)

			router.ServeHTTP(respRecorder, req)

			assert.Equal(t, tc.want.httpCode, respRecorder.Code)
			// Here we could not compare the body itself but convert the json to an struct and comparing the struct
			// If the json is misordered it would not fail.
			assert.Equal(t, tc.want.body, respRecorder.Body.String())
		})
	}
}
