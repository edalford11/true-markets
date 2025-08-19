package api_test

import (
	"github.com/edalford11/true-markets/internal/tests"
	"net/http"
)

func (suite *ApiTestSuite) TestGetPrices() {
	endpoint := suite.ListenConfig.GetEndpointAddress("/v1/prices")
	recorder := tests.ServeGet(suite, endpoint)

	expectedResponse := map[string]string{
		"BTCUSDT": "1234.56789",
		"ETSUSDT": "9876.54321",
		"BNBBTC":  "0.001",
	}

	tests.TestResponse(suite.T(), recorder, expectedResponse, http.StatusOK)
}
