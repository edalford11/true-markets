package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	V1 string = "/v1"
)

type Server struct {
	router *gin.Engine
}

func NewServer(router *gin.Engine) *Server {
	return &Server{
		router: router,
	}
}

func (s *Server) Run() error {
	s.Routes()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%v", viper.GetString("server.port")),
		Handler:           s.router.Handler(),
		ReadTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Info().Msgf("server starting on port %s", viper.GetString("server.port"))

	var cancel context.CancelFunc

	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Debug().Msg("listen and serve closed")
		} else {
			panic(err)
		}
	}()

	log.Debug().Msg("waiting for quit channel")
	<-quit

	log.Debug().Msg("quit channel received signal")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// wait 30 seconds to shutdown before timeout
	err := server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("server shutdown error")
		return err
	}

	log.Debug().Msg("server shutdown")
	return nil
}

func (s *Server) Routes() {
	s.router.GET("/health", s.HealthChecker())

	v1 := s.router.Group(V1)
	v1.GET("/price", s.getPrice)
	v1.GET("/prices", s.getPrices)
}

func (s *Server) HealthChecker() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
	}
}
