package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) GetOrderInfo(c *gin.Context) {
	orderUid := c.Param("orderUid")

	orderData, err := h.GetFromCache(orderUid)
	if err != nil {
		fmt.Println("error:", err.Error())
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusOK, orderData)
}

func (h *handler) GetAllOrders(c *gin.Context) {
	orderData := h.GetAllFromCache()
	c.JSON(http.StatusOK, orderData)
}
