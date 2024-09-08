package services

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

			// 2nd and 3rd parameters are not used in this test
			svc := NewCartItemsService(m.repo, "http://dummyhost.com", 5*time.Second)

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

// Make sure that the reservation call is correct and has all the needed params
func checkReservationHttpRequest(t *testing.T, r *http.Request) {
	assert.Equal(t, r.Method, http.MethodPost)
	assert.Equal(t, r.URL.Path, "/reserve")
	assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
	body, err := io.ReadAll(r.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"version":"1.0.0","item":{"id":"1","name":"potato","quantity":1,"reservationId":""}}`, string(body))
}

func Test_AddItemsService_GivenCartItemsServiceCreated(t *testing.T) {
	randomCartItem := model.CartItem{
		Id:       "1",
		Name:     "potato",
		Quantity: 1,
	}
	randomError := errors.New("random error")
	ctx := context.Background()

	type input struct {
		expectReservationCall bool
		timeout               time.Duration
		item                  model.CartItem
	}
	type want struct {
		err   error
		items []model.CartItem
	}
	tests := []struct {
		name                       string
		in                         input
		mocks                      func(m cartItemsServiceMocks, c chan<- int)
		reservationFakeHttpHandler func(t *testing.T, w http.ResponseWriter, r *http.Request, c chan<- int)
		want                       want
	}{
		{
			name: "WhenAddItemToCartAndFailsToAddToCart_ThenError",
			in: input{
				expectReservationCall: false,
				timeout: 40 * time.Second,
				item:    randomCartItem,
			},
			mocks: func(m cartItemsServiceMocks, c chan<- int) {
				m.repo.EXPECT().Add(gomock.Any(), randomCartItem).
					Return(randomError)
			},
			want: want{
				err: randomError,
			},
		}, {
			name: "WhenAddItemToCartAndReserveFails_ThenError",
			in: input{
				expectReservationCall: true,
				timeout: 40 * time.Second,
				item:    randomCartItem,
			},
			mocks: func(m cartItemsServiceMocks, c chan<- int) {
				m.repo.EXPECT().Add(gomock.Any(), randomCartItem).
					Return(nil)
				m.repo.EXPECT().SetReservationId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			reservationFakeHttpHandler: func(t *testing.T, w http.ResponseWriter, r *http.Request, c chan<- int) {
				checkReservationHttpRequest(t, r)
				w.WriteHeader(http.StatusNotFound) // random error
				c <- 1
			},
			want: want{
				err: nil,
			},
		}, {
			name: "WhenAddItemToCartAndTimesOut_ThenError",
			in: input{
				expectReservationCall: true,
				timeout: 1 * time.Second,
				item:    randomCartItem,
			},
			mocks: func(m cartItemsServiceMocks, c chan<- int) {
				m.repo.EXPECT().Add(gomock.Any(), randomCartItem).
					Return(nil)
				m.repo.EXPECT().SetReservationId(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			reservationFakeHttpHandler: func(t *testing.T, w http.ResponseWriter, r *http.Request, c chan<- int) {
				checkReservationHttpRequest(t, r)
				
				// the request will time out and the SetReservationId method will not be called
				time.Sleep(4 * time.Second)
				c <- 1
			},
			want: want{
				err: nil,
			},
		}, {
			name: "WhenAddItemToCartAndOK_ThenReservationIsWrittenInRepo",
			in: input{
				expectReservationCall: true,
				timeout: 3 * time.Second,
				item:    randomCartItem,
			},
			mocks: func(m cartItemsServiceMocks, c chan<- int) {
				m.repo.EXPECT().Add(gomock.Any(), randomCartItem).
					Return(nil)
				m.repo.EXPECT().SetReservationId(gomock.Any(), gomock.Any(), "fancyReservationId").
					DoAndReturn(func(ctx context.Context, item model.CartItem, reservationId string) error {
						c<-1 // channel needs to be fed when setReservationId is called
						return nil
					})
			},
			reservationFakeHttpHandler: func(t *testing.T, w http.ResponseWriter, r *http.Request, c chan<- int) {
				checkReservationHttpRequest(t, r)
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK) // random error
				w.Write([]byte(`{"version":"1.0.0","reservationId":"fancyReservationId"}`))
							
			},
			want: want{
				err: nil,
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

			done := make(chan int)
			tc.mocks(m, done)

			// http test to mock reservation server
			reservationFakeHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				tc.reservationFakeHttpHandler(t, w, r, done)
			}))
			defer reservationFakeHttpServer.Close()
			reservationHost := reservationFakeHttpServer.URL

			// 2nd and 3rd parameters are not used in this test
			svc := NewCartItemsService(m.repo, reservationHost, tc.in.timeout)

			addErr := svc.Add(ctx, tc.in.item)
			if tc.want.err != nil {
				// helps being agnostic with the error message, as it can change and wrongfully break the tests
				// If an error comprobation is needed, assert.ErrorIs or ErrorAs can be used.
				assert.Error(t, addErr)
			} else {
				assert.NoError(t, addErr)
			}

			// we check with the channel that the http call has been made. 
			// As it is run in a goroutine, most probably, the funcion will return before the call is even made
			if tc.in.expectReservationCall {
				select {
				case <-done:
					// Successfully received the http call to the reserve endpoint
				case <-time.After(10 * time.Second):
					t.Error("Timed out waiting for goroutine to finish")
				}
			}

		})
	}
}
