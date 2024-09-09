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

// Identical to Cart Item Request, but since they are 2 diffent usages, they could diverge from each other.
// Reusing the same struct could be error prone and difficult to mantain, if some field is added for one case,
// but not needed for the other one
type ItemReservationRequest struct {
	Version string   `json:"version"`
	Item    CartItem `json:"item"`
}

func NewItemReservationRequest(item CartItem) ItemReservationRequest {
	return ItemReservationRequest{
		Version: "1.0.0",
		Item:    item,
	}
}

type ItemReservationResponse struct {
	Version       string `json:"version"`
	ReservationId string `json:"reservationId"`
}
