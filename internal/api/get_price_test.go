package api_test

import (
	"github.com/edalford11/true-markets/internal/tests"
	"net/http"
)

func (suite *ApiTestSuite) TestGetPrice() {
	endpoint := suite.ListenConfig.GetEndpointAddress("/v1/price?symbols=btcusdt,etsusdt,invalid")
	recorder := tests.ServeGet(suite, endpoint)

	expectedResponse := map[string]string{
		"BTCUSDT": "1234.56789",
		"ETSUSDT": "9876.54321",
	}

	tests.TestResponse(suite.T(), recorder, expectedResponse, http.StatusOK)

	endpoint = suite.ListenConfig.GetEndpointAddress("/v1/price?symbols=BTCUSDT,ETSUSDT,INVALID")
	recorder = tests.ServeGet(suite, endpoint)

	tests.TestResponse(suite.T(), recorder, expectedResponse, http.StatusOK)
}
