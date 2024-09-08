package model

type CartItem struct {
	Name          string `json:"name"`
	Quantity      int    `json:"quantiy"`
	ReservationId string `json:"reservationId,omitemtpy"`
}
