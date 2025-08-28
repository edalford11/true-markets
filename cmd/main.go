package main

import (
	"fmt"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/edalford11/true-markets/config"
	"github.com/edalford11/true-markets/internal/api"
	"github.com/edalford11/true-markets/internal/helpers"
	"github.com/edalford11/true-markets/internal/util"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func main() {
	config.Init()

	cmdAPI := &cobra.Command{
		Use:   "api",
		Short: "Run the TrueMarkets API",
		Long:  `Run an instance of the True Markets API`,
		Run: func(cmd *cobra.Command, args []string) {
			err := RunAPI()
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd := &cobra.Command{Use: "truemarkets"}
	rootCmd.AddCommand(cmdAPI)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func RunAPI() error {
	util.InitLogger()

	stopCh := connectToBinance()

	environment := config.Environment(viper.GetString("environment"))

	if environment == config.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	// skip logging for health check
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/health"))

	if environment == config.Production || environment == config.Staging {
		helpers.InitSentry(environment)
		defer sentry.Flush(2 * time.Second)

		r.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
	}

	s := api.NewServer(r)
	if err := s.Run(); err != nil {
		return err
	}

	fmt.Println("Received shutdown signal, stopping binance listeners")
	stopCh <- struct{}{}

	return nil
}

func connectToBinance() chan struct{} {
	symbols := strings.Split(viper.GetString("binance.symbols"), ",")
	log.Info().Msgf("Connecting to binance with symbols: %s", strings.Join(symbols, ","))

	mainStopCh := make(chan struct{})

	go func() {
		for {
			select {
			case <-mainStopCh:
				return
			default:
				log.Info().Msg("Establishing Binance websocket connection")

				websocketStreamClient := binance_connector.NewWebsocketStreamClient(true, "wss://stream.binance.us:9443")

				wsHandler := func(event *binance_connector.WsCombinedTradeEvent) {
					symbol := event.Data.Symbol
					price := event.Data.Price
					util.GetSymbolPriceMap().Set(strings.ToUpper(symbol), price)
				}

				reconnectCh := make(chan struct{})

				errHandler := func(err error) {
					log.Error().Err(err).Msg("Binance websocket error - triggering reconnect")
					reconnectCh <- struct{}{}
				}

				_, stopCh, err := websocketStreamClient.WsCombinedTradeServe(symbols, wsHandler, errHandler)
				if err != nil {
					log.Error().Err(err).Msg("Failed to connect to Binance websocket, retrying in 5 seconds")
					time.Sleep(5 * time.Second)
					continue
				}

				log.Info().Msg("Binance websocket connection established")

				select {
				case <-mainStopCh:
					stopCh <- struct{}{}
					return
				case <-reconnectCh:
					log.Info().Msg("Reconnecting to Binance websocket in 5 seconds")
					time.Sleep(5 * time.Second)
				}
			}
		}
	}()

	return mainStopCh
}
