package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polonkoevv/wb-tech/internal/service"
)

type handler struct {
	service.ServiceInterface
}

func New(serv service.ServiceInterface) *gin.Engine {
	h := handler{
		serv,
	}

	r := gin.New()

	r.Static("/css", "templates/css/")
	r.LoadHTMLGlob("templates/index.html")

	r.GET("/order/:orderUid", h.GetOrderInfo)
	r.GET("/order/", h.GetAllOrders)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	return r
}
