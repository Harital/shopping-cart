package model

type CartItem struct {
	Name          string `json:"name"`
	Quantity      int    `json:"quantiy"`
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