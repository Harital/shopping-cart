package http

import (
	"net/http"

	"github.com/Harital/shopping-cart/internal/core/model"
	"github.com/Harital/shopping-cart/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CartItemshandler struct {
	router          *gin.RouterGroup
	cartItemService ports.CartItemsService
}

func NewCartItemsHandler(r *gin.RouterGroup, svc ports.CartItemsService) *CartItemshandler {
	return &CartItemshandler{
		router:          r,
		cartItemService: svc,
	}
}

// This functionality could be done in the constructor, but it is not a good practise to add
// logic (that can potentially fail) in the constructor
func (cih *CartItemshandler) Register() {
	cih.router.GET("/items", cih.getCartItems)
	cih.router.POST("/items", cih.addCartItems)
}

func (cih CartItemshandler) getCartItems(c *gin.Context) {
	items, getErr := cih.cartItemService.Get(c)

	// Other error types should be checked here, like if the user is properly authenticated.
	// The response should be different dependint on the error type
	if getErr != nil {
		resp := model.NewErrorResponse("Error getting the cart items")
		log.
			Error().
			Err(getErr).
			Msg("when getting cart items")
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp := model.NewGetCartITemsResponse(&items)
	c.JSON(http.StatusOK, *resp)
}

func (cih CartItemshandler) addCartItems(c *gin.Context) {
	var item model.CartItemRequest
	if bindErr := c.ShouldBindJSON(&item); bindErr != nil {
		resp := model.NewErrorResponse("bad request")
		log.
			Error().
			Err(bindErr).
			Msg("Bad json request")
		// We could consider here to print the body that caused this error for debugbging purposes
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// We could make more input checks, like all the require fields are filled, they have the proper format, etc

	addErr := cih.cartItemService.Add(c, item.Item)
	if addErr != nil {
		resp := model.NewErrorResponse("internal error")
		log.
			Error().
			Err(addErr).
			Msg("adding an item to the basket")
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.Status(http.StatusAccepted)
}
