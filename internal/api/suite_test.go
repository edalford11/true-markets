package api_test

import (
	"fmt"
	"github.com/edalford11/true-markets/internal/util"
	"testing"

	"github.com/edalford11/true-markets/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ListeningConfig struct {
	host string
	port string
}

// GetEndpointAddress constructs the endpoint address from viper env vars.
// must start with a forward slash.
func (l ListeningConfig) GetEndpointAddress(endpointPath string, any ...any) string {
	return fmt.Sprintf(fmt.Sprintf("http://%s:%s%s", l.host, l.port, endpointPath), any...)
}

type ApiTestSuite struct {
	suite.Suite
	router       *gin.Engine
	ListenConfig ListeningConfig
}

func (suite *ApiTestSuite) Router() *gin.Engine {
	return suite.router
}

func (suite *ApiTestSuite) SetupSuite() {
	suite.ListenConfig = ListeningConfig{
		host: "localhost",
		port: "8080",
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	s := api.NewServer(r)
	s.Routes()
	suite.router = r
}

func (suite *ApiTestSuite) BeforeTest(suiteName, testName string) {
	util.GetSymbolPriceMap().DeleteAll()
	util.GetSymbolPriceMap().Set("BTCUSDT", "1234.56789")
	util.GetSymbolPriceMap().Set("ETSUSDT", "9876.54321")
	util.GetSymbolPriceMap().Set("BNBBTC", "0.001")
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
