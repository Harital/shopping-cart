package model

type CartItem struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Quantity      int    `json:"quantity"`
	ReservationId string `json:"reservationId,omitemtpy"`
}

type GetCartItemsResponse struct {
	Version string     `json:"Version"`
	Items   []CartItem `json:"Items"`
}

func NewGetCartITemsResponse(items *[]CartItem) *GetCartItemsResponse {
	return &GetCartItemsResponse{
		Version: "1.0.0",
		// we could set the GetCartItemsResponse.Items as a pointer, but we would need to allocate space
		Items: *items,
	}
}

type CartItemRequest struct {
	Version string   `json:"version"`
	Item    CartItem `json:"item"`
}
