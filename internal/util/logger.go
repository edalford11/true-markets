package util

import (
	"os"
	"time"

	"github.com/edalford11/true-markets/config"

	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	if config.IsLocalEnvironment() {
		log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	}

	logLevel := viper.GetString("logger.level")
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
}

type ginRequest struct {
	SerName    string
	Path       string
	Latency    time.Duration
	Method     string
	StatusCode int
	ClientIP   string
	MsgStr     string
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		InitLogger() // Reset logger context vars

		t := time.Now()
		// before request
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		// after request
		// latency := time.Since(t)
		// clientIP := c.ClientIP()
		// method := c.Request.Method
		// statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		msg := c.Errors.String()
		if msg == "" {
			msg = "Request"
		}

		cData := &ginRequest{
			Path:       path,
			Latency:    time.Since(t),
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			ClientIP:   c.ClientIP(),
			MsgStr:     msg,
		}

		logSwitch(cData)
	}
}

func logSwitch(data *ginRequest) {
	switch {
	case data.StatusCode >= 400 && data.StatusCode < 500:
		{
			log.Warn().Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
		}
	case data.StatusCode >= 500:
		{
			log.Error().Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
		}
	default:
		log.Info().Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
	}
}
