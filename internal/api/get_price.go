package api

import (
	"github.com/edalford11/true-markets/internal/util"
	api "github.com/edalford11/true-markets/pkg/models"
	"github.com/edalford11/true-markets/pkg/models/requests"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func (s *Server) getPrice(ctx *gin.Context) {
	queryParams := new(requests.PriceQueryParams)
	err := ctx.ShouldBindQuery(queryParams)
	if err != nil {
		log.Logger.Error().Msgf("invalid url params: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &api.ErrorResponse{Error: err.Error()})
		return
	}

	symbolPrices := make(map[string]string, len(queryParams.Symbols))

	for _, symbol := range strings.Split(queryParams.Symbols, ",") {
		price, ok := util.GetSymbolPriceMap().Get(strings.ToUpper(symbol))
		if ok {
			symbolPrices[strings.ToUpper(symbol)] = price
		}
	}

	ctx.JSON(http.StatusOK, symbolPrices)
}
