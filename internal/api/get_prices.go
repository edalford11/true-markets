package api

import (
	"github.com/edalford11/true-markets/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) getPrices(ctx *gin.Context) {
	priceMap := util.GetSymbolPriceMap().GetAll()
	ctx.JSON(http.StatusOK, priceMap)
}
