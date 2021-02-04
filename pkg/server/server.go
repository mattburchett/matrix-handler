package server

import (
	"fmt"
	"net/http"

	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/config"
	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/generic"

	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/router"
	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/slack"

	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Run ...
func Run(info *router.BuildInfo) error {
	conf, err := config.GetConfig("config.json")
	if err != nil {
		log.Fatal().Err(err)
	}

	level, err := zerolog.ParseLevel(conf.LogLevel)
	if err != nil {
		level = zerolog.ErrorLevel
		log.Warn().Err(err).Msgf("unable to parse log level, logging level is set to %s", level.String())
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = log.With().Str("app", conf.Name).Logger()

	router := router.NewRouter(info)

	router.HandleWithMetrics("/generic/{matrixRoom}/{matrixUser}/{matrixPassword}", generic.Handle(conf)).Methods(http.MethodPost)
	router.HandleWithMetrics("/slack/{matrixRoom}/{matrixUser}/{matrixPassword}", slack.Handle(conf)).Methods(http.MethodPost)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: cors.Default().Handler(router),
	}

	log.Info().Msgf("Server running on %v", srv.Addr)
	return srv.ListenAndServe()
}
